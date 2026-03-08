#!/bin/sh
# awsexec: execute command with aws configure export-credentials
#
# Usage:
#   awsexec profile -- command [args...]
#   AWS_PROFILE=profile awsexec -- command [args...]

set -e

# Find the -- separator
separator_index=-1
i=1
while [ "${i}" -le "${#}" ]; do
    # Use eval to get argument by index (POSIX sh doesn't support ${!i})
    eval "arg=\$${i}"
    if [ "${arg}" = "--" ]; then
        separator_index="${i}"
        break
    fi
    i=$((i + 1))
done

if [ "${separator_index}" -eq -1 ]; then
    echo "Error: missing '--' separator" >&2
    exit 1
fi

# Extract profile from arguments or environment variable
if [ "${separator_index}" -gt 1 ]; then
    # Profile is specified as an argument
    profile="${1}"
else
    # Profile should be in AWS_PROFILE environment variable
    profile="${AWS_PROFILE:-}"
    if [ -z "${profile}" ]; then
        echo "Error: AWS_PROFILE environment variable is required when profile is not specified as an argument" >&2
        exit 1
    fi
fi

# Extract command after --
# Note: Variables in arithmetic expressions $((...)) don't need quotes
if [ $((${separator_index} + 1)) -gt "${#}" ]; then
    echo "Error: no command specified after '--'" >&2
    exit 1
fi

# Build command arguments by shifting to the position after --
# We need to shift separator_index times to get to the position after --
shift_count="${separator_index}"
while [ "${shift_count}" -gt 0 ]; do
    shift
    # Note: Variables in arithmetic expressions $((...)) don't need quotes
    shift_count=$((${shift_count} - 1))
done
# Now ${@} contains only the command and its arguments

# Resolve region: profile config, then AWS_DEFAULT_REGION, then AWS_REGION
region=""
if [ -n "${profile}" ]; then
    region=$(aws configure get region --profile "${profile}" 2>/dev/null) || true
else
    region=$(aws configure get region 2>/dev/null) || true
fi
[ -z "${region}" ] && region="${AWS_DEFAULT_REGION:-}"
[ -z "${region}" ] && region="${AWS_REGION:-}"

# Export credentials using aws configure export-credentials
# The output includes AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN, and AWS_CREDENTIAL_EXPIRATION
if [ -n "${profile}" ]; then
    eval "$(aws configure export-credentials --format env --profile "${profile}")"
else
    eval "$(aws configure export-credentials --format env)"
fi

# Reset AWS_PROFILE so the subprocess uses only the credential env vars
unset AWS_PROFILE

# Set region for the subprocess (same order as Go implementation)
[ -n "${region}" ] && export AWS_DEFAULT_REGION="${region}" && export AWS_REGION="${region}"

exec "${@}"
