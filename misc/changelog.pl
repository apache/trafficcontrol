#!/usr/bin/perl

use strict;
use warnings;

use WWW::Curl::Easy;
use JSON;

my $owner = shift;
my $repo = shift;
my $milestone = shift;
my $url = "https://api.github.com";

sub milestone_lookup
{
  my $url = shift;
  my $owner = shift;
  my $repo = shift;
  my $milestone_title = shift;
  my $endpoint = "/repos/$owner/$repo/milestones";

  my $params = "state=all";

  my $resp_body;
  my $curl = WWW::Curl::Easy->new;

  #$curl->setopt(CURLOPT_VERBOSE, 1);
  $curl->setopt(CURLOPT_HTTPHEADER, ['Accept: application/vnd.github.v3+json', 'User-Agent: Awesome-Octocat-App']);
  $curl->setopt(CURLOPT_WRITEDATA, \$resp_body);
  $curl->setopt(CURLOPT_URL, $url . $endpoint . '?' . $params);

  my $retcode = $curl->perform();
  if ($retcode == 0 && $curl->getinfo(CURLINFO_HTTP_CODE) == 200)
  {
    my $milestones = from_json($resp_body);
    foreach my $milestone (@{ $milestones })
    {
      if ($milestone->{title} eq $milestone_title)
      {
        return $milestone->{number};
      }
    }
  }

  return undef;
}

sub issue_search
{
  my $url = shift;
  my $owner = shift;
  my $repo = shift;
  my $milestone_id = shift;
  my $page = shift;
  my $endpoint = "/repos/$owner/$repo/issues";

  my $params = "milestone=$milestone_id&state=closed&page=$page";

  my $resp_body;
  my $curl = WWW::Curl::Easy->new;

  #$curl->setopt(CURLOPT_VERBOSE, 1);
  $curl->setopt(CURLOPT_HTTPHEADER, ['Accept: application/vnd.github.v3+json', 'User-Agent: Awesome-Octocat-App']);
  $curl->setopt(CURLOPT_WRITEDATA, \$resp_body);
  $curl->setopt(CURLOPT_URL, $url . $endpoint . '?' . $params);

  my $retcode = $curl->perform();
  if ($retcode == 0 && $curl->getinfo(CURLINFO_HTTP_CODE) == 200) {
    return from_json($resp_body);
  }

  undef;
}

my $milestone_id = milestone_lookup($url, $owner, $repo, $milestone);

if (!defined($milestone_id))
{
  exit 1;
}

my $issues;
my $changelog;
my $page = 1;

do {
  $issues = issue_search($url, $owner, $repo, $milestone_id, $page);
  foreach my $issue (@{ $issues })
  {
    if (defined($issue))
    {
      push @{ $changelog }, {number => $issue->{number},  title => $issue->{title}};
    }
  }
  $page++;
} while (scalar @{ $issues });

if (defined($changelog))
{
  print "Changes with Traffic Control $milestone\n";

  foreach my $issue (sort {$a->{number} <=> $b->{number}} @{ $changelog })
  {
    print "  #$issue->{number} - $issue->{title}\n";
  }
}
