package Term::ReadPassword;

use strict;
use Term::ReadLine;
use POSIX qw(:termios_h);
my %CC_FIELDS = (
	VEOF => VEOF,
	VEOL => VEOL,
	VERASE => VERASE,
	VINTR => VINTR,
	VKILL => VKILL,
	VQUIT => VQUIT,
	VSUSP => VSUSP,
	VSTART => VSTART,
	VSTOP => VSTOP,
	VMIN => VMIN,
	VTIME => VTIME,
    );

use vars qw(
    $VERSION @ISA @EXPORT @EXPORT_OK
    $ALLOW_STDIN %SPECIAL $SUPPRESS_NEWLINE $INPUT_LIMIT
    $USE_STARS $STAR_STRING $UNSTAR_STRING
);

require Exporter;

@ISA = qw(Exporter);
@EXPORT = qw(
	read_password 
);
$VERSION = '0.11';

# The special characters in the input stream
%SPECIAL = (
    "\x03"	=> 'INT',	# Control-C, Interrupt
    "\x15"	=> 'NAK',	# Control-U, NAK (clear buffer)
    "\x08"	=> 'DEL',	# Backspace
    "\x7f"	=> 'DEL',	# Delete
    "\x0d"	=> 'ENT',	# CR, Enter
    "\x0a"	=> 'ENT',	# LF, Enter
);

# The maximum amount of data for the input buffer to hold
$INPUT_LIMIT = 1000;

sub read_password {
    my($prompt, $idle_limit, $interruptable) = @_;
    $prompt = '' unless defined $prompt;
    $idle_limit = 0 unless defined $idle_limit;
    $interruptable = 0 unless defined $interruptable;

    # Let's open the TTY (rather than STDIN) if we can
    local(*TTY, *TTYOUT);
    my($in, $out) = Term::ReadLine->findConsole;
    die "No console available" unless $in;
    if (open TTY, "+<$in") {
        # Cool
    } elsif ($ALLOW_STDIN) {
        open TTY, "<&STDIN"
	    or die "Can't re-open STDIN: $!";
    } else {
        die "Can't open '$in' read/write: $!";
    }

    # And let's send the output to the TTY as well
    if (open TTYOUT, ">>$out") {
	# Cool
    } elsif ($ALLOW_STDIN) {
	# Well, let's allow STDOUT as well
	open TTYOUT, ">>&STDOUT"
	    or die "Can't re-open STDOUT: $!";
    } else {
	die "Can't open '$out' for output: $!";
    }

    # Don't buffer it!
    select( (select(TTYOUT), $|=1)[0] );
    print TTYOUT $prompt;

    # Okay, now remember where everything was, so we can put it back when
    # we're done 
    my $fd_tty = fileno(TTY);
    my $term = POSIX::Termios->new();
    $term->getattr($fd_tty);
    my $original_flags = $term->getlflag();
    my %original_cc;
    for my $field_name (keys %CC_FIELDS) {
        $original_cc{$field_name} = $term->getcc($CC_FIELDS{$field_name});
    }

    # What makes this setup different from the ordinary?
    # No keyboard-generated signals, no echoing, no canonical input
    # processing (like backspace handling)
    my $flags = $original_flags & ~(ISIG | ECHO | ICANON);
    $term->setlflag($flags);
    if ($idle_limit) {
	# $idle_limit is in seconds, so multiply by ten
	$term->setcc(VTIME, 10 * $idle_limit);
	# Continue running the program after that time, even if there
	# weren't any characters typed
	$term->setcc(VMIN, 0);
    } else {
	# No time limit, but...
	$term->setcc(VTIME, 0);
	# Continue as soon as one character has been struck
	$term->setcc(VMIN, 1);
    }

    # Optionally echo stars in place of password characters. The 
    # $unstar_string uses backspace characters.
    my $star_string = $USE_STARS ? ($STAR_STRING || '*') : '';
    my $unstar_string = $USE_STARS ? ($UNSTAR_STRING || "\b*\b \b") : '';

    # If there's anything already buffered, we should throw it out. This
    # is to discourage users from typing their password before they see
    # the prompt, since their keystrokes may be echoing on the screen. 
    #
    # So this statement supposedly makes sure the prompt goes out, the
    # unread input buffer is discarded, and _then_ the changes take
    # effect. Thus, everything they typed ahead is (probably) echoed.
    $term->setattr($fd_tty, TCSAFLUSH);

    my $input = '';
    my $return_value;
KEYSTROKE:
    while (1) {
        my $new_keys = '';
	my $count = sysread(TTY, $new_keys, 99);
	# We're here, so either the idle_limit expired, or the user typed
	# something.
	if ($count) {
	    for my $new_key (split //, $new_keys) {
	        if (my $meaning = $SPECIAL{$new_key}) {
		    if ($meaning eq 'ENT') {
		        # Enter/return key
			# Return what we have so far
			$return_value = $input;
			last KEYSTROKE;
		    } elsif ($meaning eq 'DEL') {
		        # Delete/backspace key
			# Take back one char, if possible
			if (length $input) {
			    $input = substr $input, 0, length($input)-1;
			    print TTYOUT $unstar_string;
			}
		    } elsif ($meaning eq 'NAK') {
		        # Control-U (NAK)
		        # Clear what we have read so far
			for (1..length $input) {
			    print TTYOUT $unstar_string;
			}
		        $input = '';
		    } elsif ($interruptable and $meaning eq 'INT') {
			# Breaking out of the program
			# Return early
			last KEYSTROKE;
		    } else {
		        # Just an ordinary keystroke
			$input .= $new_key;
			print TTYOUT $star_string;
		    }
		} else {
		    # Not special
		    $input .= $new_key;
		    print TTYOUT $star_string;
		}
 	    }
	    # Just in case someone sends a lot of data
	    $input = substr($input, 0, $INPUT_LIMIT)
	        if length($input) > $INPUT_LIMIT;
	} else {
	    # No count, so something went wrong. Assume timeout.
	    # Return early
	    last KEYSTROKE;
	}
    }

    # Done with waiting for input. Let's not leave the cursor sitting
    # there, after the prompt.
    print TTYOUT "\n" unless $SUPPRESS_NEWLINE;

    # Let's put everything back where we found it.
    $term->setlflag($original_flags);
    while (my($field, $value) = each %original_cc) {
        $term->setcc($CC_FIELDS{$field}, $value);
    }
    $term->setattr($fd_tty, TCSAFLUSH);
    close(TTY);
    close(TTYOUT);
    $return_value;
}

1;
__END__

=head1 NAME

Term::ReadPassword - Asking the user for a password

=head1 SYNOPSIS

  use Term::ReadPassword;
  while (1) {
    my $password = read_password('password: ');
    redo unless defined $password;
    if ($password eq 'flubber') {
      print "Access granted.\n";
      last;
    } else {
      print "Access denied.\n";
      redo;
    }
  }

=head1 DESCRIPTION

This module lets you ask the user for a password in the traditional way,
from the keyboard, without echoing.

This is not intended for use over the web; user authentication over the
web is another matter entirely. Also, this module should generally be used
in conjunction with Perl's B<crypt()> function, sold separately.

The B<read_password> function prompts for input, reads a line of text from
the keyboard, then returns that line to the caller. The line of text
doesn't include the newline character, so there's no need to use B<chomp>.

While the user is entering the text, a few special characters are processed.
The character delete (or the character backspace) will back up one
character, removing the last character in the input buffer (if any). The
character CR (or the character LF) will signal the end of input, causing the
accumulated input buffer to be returned. Control-U will empty the input
buffer. And, optionally, the character Control-C may be used to terminate
the input operation. (See details below.) All other characters, even ones
which would normally have special purposes, will be added to the input
buffer.

It is not recommended, though, that you use the as-yet-unspecified control
characters in your passwords, as those characters may become meaningful in
a future version of this module. Applications which allow the user to set
their own passwords may wish to enforce this rule, perhaps with code
something like this:

    {
      # Naked block for scoping and redo
      my $new_pw = read_password("Enter your new password: ");
      if ($new_pw =~ /([^\x20-\x7E])/) {
        my $bad = unpack "H*", $1;
	print "Your password may not contain the ";
	print "character with hex code $bad.\n";
	redo;
      } elsif (length($new_pw) < 5) {
        print "Your password must be longer than that!\n";
	redo;
      } elsif ($new_pw ne read_password("Enter it again: ")) {
	print "Passwords don't match.\n";
	redo;
      } else {
        &change_password($new_pw);
	print "Your password is now changed.\n";
      }
    }

The second parameter to B<read_password> is the optional C<idle_timeout>
value. If it is a non-zero number and there is no keyboard input for that
many seconds, the input operation will terminate. Notice that this is not
an overall time limit, as the timer is restarted with each new character.

The third parameter will optionally allow the input operation to be
terminated by the user with Control-C. If this is not supplied, or is
false, a typed Control-C will be entered into the input buffer just as any
other character. In that case, there is no way from the keyboard to
terminate the program while it is waiting for input. (That is to say, the
normal ability to generate signals from the keyboard is suspended during
the call to B<read_password>.)

If the input operation terminates early (either because the idle_timeout
was exceeded, or because a Control-C was enabled and typed), the return
value will be C<undef>. In either case, there is no way provided to
discover what (if anything) was typed before the early termination, or why
the input operation was terminated.

So as to discourage users from typing their passwords anywhere except at
the prompt, any input which has been "typed ahead" before the prompt
appears will be discarded. And whether the input operation terminates
normally or not, a newline character will be printed, so that the cursor
will not remain on the line after the prompt. 

=head1 BUGS

Windows users will want Term::ReadPassword::Win32.

This module has a poorly-designed interface, and should be thoroughly
rethought and probably re-united with the Windows version. 

Users who wish to see password characters echoed as stars may set
$Term::ReadPassword::USE_STARS to a true value. The bugs are that some
terminals may not erase stars when the user corrects an error, and that
using stars leaks information to shoulder-surfers. 

=head1 SECURITY

You would think that a module dealing with passwords would be full of
security features. You'd think that, but you'd be wrong. For example, perl
provides no way to erase a piece of data from memory. (It's easy to erase
it so that it can't be accessed from perl, but that's not the same thing
as expunging it from the actual memory.) If you've entered a password,
even if the variable that contained that password has been erased, it may
be possible for someone to find that password, in plaintext, in a core
dump. And that's just one potential security hole.

In short, if serious security is an issue, don't use this module.

=head1 LICENSE

This program is free software; you may redistribute it, modify it, or
both, under the same terms as Perl itself.

=head1 AUTHOR

Tom Phoenix <rootbeer@redcat.com>. Copyright (C) 2007 Tom Phoenix.

=head1 SEE ALSO

Term::ReadLine, L<perlfunc/crypt>, and your system's manpages for the
low-level I/O operations used here.

=cut
