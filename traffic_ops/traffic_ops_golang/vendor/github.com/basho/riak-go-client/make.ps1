<#
.SYNOPSIS
    Powershell script to build Riak Go Client on Windows
.DESCRIPTION
    This script ensures that your build environment is sane and will run 'go' correctly depending on parameters passed to this script.
.PARAMETER Target
    Target to build. Can be one of the following:
        * ProtoGen           - generate *.pb.go files from *.proto files.
                               Requires Go protoc utility (https://github.com/golang/protobuf)
        * Format             - run *.go files through 'go fmt'
        * Test               - Run 'go vet' and all tests (default target)
        * UnitTest           - Run unit tests
        * IntegrationTest    - Run live integration tests
        * IntegrationTestHll - Run Hyperloglog integration tests
        * TimeseriesTest     - Run live timeseries tests
.PARAMETER Verbose
    Use to increase verbosity.
.EXAMPLE
    C:\Users\Bashoman> cd go\src\github.com\basho\riak-go-client
    C:\Users\Bashoman\go\src\github.com\basho\riak-go-client>.\make.ps1 -Target Protoc -Verbose
.NOTES
    Author: Luke Bakken
    Date:   June 1, 2015
#>
[CmdletBinding()]
Param(
    [Parameter(Mandatory=$False, Position=0)]
    [ValidateSet('ProtoGen', 'Format',
        'Test', 'UnitTest', 'IntegrationTest', 'IntegrationTestHll', 'TimeseriesTest',
        IgnoreCase = $True)]
    [string]$Target = 'Test'
)

Set-StrictMode -Version Latest

$package = 'github.com/basho/riak-go-client'

$IsDebug = $DebugPreference -ne 'SilentlyContinue'
$IsVerbose = $VerbosePreference -ne 'SilentlyContinue'

# Note:
# Set to Continue to see DEBUG messages
if ($IsVerbose) {
    $DebugPreference = 'Continue'
}

trap
{
    Write-Error -ErrorRecord $_
    exit 1
}

function Get-ScriptPath {
    $scriptDir = Get-Variable PSScriptRoot -ErrorAction SilentlyContinue | ForEach-Object { $_.Value }
    if (!$scriptDir) {
        if ($MyInvocation.MyCommand.Path) {
            $scriptDir = Split-Path $MyInvocation.MyCommand.Path -Parent
        }
    }
    if (!$scriptDir) {
        if ($ExecutionContext.SessionState.Module.Path) {
            $scriptDir = Split-Path (Split-Path $ExecutionContext.SessionState.Module.Path)
        }
    }
    if (!$scriptDir) {
        $scriptDir = $PWD
    }
    return $scriptDir
}

function Do-ProtoGen {
    $script_path = Get-ScriptPath
    $rpb_path = Join-Path -Path $script_path -ChildPath 'rpb'
    $proto_path = Join-Path -Path $script_path -ChildPath (Join-Path -Path 'riak_pb' -ChildPath 'src')
    $proto_wild = Join-Path -Path $proto_path -ChildPath '*.proto'
    Get-ChildItem $proto_wild | ForEach-Object {
        $proto_file_basename = $_.BaseName
        $rpb_path_tmp = Join-Path -Path $rpb_path -ChildPath $proto_file_basename
        If (!(Test-Path $rpb_path_tmp)) {
            New-Item $rpb_path_tmp -Type Directory -Force
        }
        Write-Verbose "protoc: --go_out=$rpb_path_tmp --proto_path=$proto_path $_"
        & { protoc --go_out=$rpb_path_tmp --proto_path=$proto_path $_ }
        if ($? -ne $True) {
            throw "protoc.exe failed: $LastExitCode"
        }

        $rpb_file = Join-Path -Path $rpb_path_tmp -ChildPath "$proto_file_basename.pb.go"
        Write-Verbose "post-processing $rpb_file"

        (Get-Content $rpb_file) |
            ForEach-Object {
                $_ -Replace 'import proto "code.google.com/p/goprotobuf/proto"', 'import proto "github.com/golang/protobuf/proto"'
            } | Set-Content $rpb_file

        if ($_.Name -eq 'riak_search.proto' -or $_.Name -eq 'riak_kv.proto') {
            (Get-Content $rpb_file) |
                ForEach-Object {
                    $_ -Replace 'import riak "riak.pb"', 'import riak "github.com/basho/riak-go-client/rpb/riak"'
                } | Set-Content $rpb_file
        }
    }
}

function Execute($cmd, $argz) {
    Write-Verbose "$cmd $argz"
    & $cmd $argz
    if ($? -ne $True) {
        throw "'$cmd $argz' failed: $LastExitCode"
    }
    Write-Debug "'$cmd $argz' exit code: $LastExitCode"
}

function Do-Format {
    $script_path = Get-ScriptPath
    Write-Verbose "go fmt $script_path"
    $cmd = 'gofmt'
    $argz = '-s', '-w', $script_path
    Execute $cmd $argz
}

function Do-Vet {
    $cmd = 'go.exe'
    $script_path = Get-ScriptPath
    $argz = 'tool', 'vet', '-shadow=true', '-shadowstrict=true', $script_path
    Execute $cmd $argz
    $argz = 'vet', $package
    Execute $cmd $argz
}

function Do-UnitTest {
    $v = ''
    if ($IsVerbose) {
        $v = '-v'
    }
    $cmd = 'go.exe'
    $argz = 'test', $v, "$package/..."
    Execute $cmd $argz
}

function Do-IntegrationTest {
    $v = ''
    if ($IsVerbose) {
        $v = '-v'
    }
    $cmd = 'go.exe'
    $argz = 'test', $v, '-tags=integration', "$package/..."
    Execute $cmd $argz
}

function Do-IntegrationTestHll {
    $v = ''
    if ($IsVerbose) {
        $v = '-v'
    }
    $cmd = 'go.exe'
    $argz = 'test', $v, '-tags=integration_hll', "$package/..."
    Execute $cmd $argz
}

function Do-TimeseriesTest {
    $v = ''
    if ($IsVerbose) {
        $v = '-v'
    }
    $cmd = 'go.exe'
    $argz = 'test', $v, '-tags=timeseries', "$package/..."
    Execute $cmd $argz
}

Write-Debug "Target: $Target"

switch ($Target)
{
    'ProtoGen' { Do-ProtoGen }
    'Format' { Do-Format }
    'Test' { Do-Vet; Do-IntegrationTest }
    'UnitTest' { Do-Vet; Do-UnitTest }
    'IntegrationTest' { Do-Vet; Do-IntegrationTest }
    'IntegrationTestHll' { Do-Vet; Do-IntegrationTestHll }
    'TimeseriesTest' { Do-Vet; Do-TimeseriesTest }
     default { throw "Unknown target: $Target" }
}

exit 0
