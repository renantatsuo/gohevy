#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if [[ -z "${HEVY_API_KEY:-}" ]]; then
	echo "HEVY_API_KEY is required" >&2
	exit 1
fi

resource_log="${HEVY_INTEGRATION_RESOURCE_LOG:-integration_api_created_resources.log}"
test_log="${HEVY_INTEGRATION_TEST_LOG:-integration_api_test.log}"

rm -f "$resource_log" "$test_log"

echo "Running live Hevy API integration tests..."
echo "Created resources will be written to: $resource_log"
echo "Full test output will be written to: $test_log"
echo

set +e
HEVY_INTEGRATION_RESOURCE_LOG="$resource_log" go test -tags=integration -v -count=1 -timeout=180s ./... 2>&1 | tee "$test_log"
status=${PIPESTATUS[0]}
set -e

echo
if [[ -s "$resource_log" ]]; then
	echo "Created resources log ($resource_log):"
	while IFS= read -r line; do
		echo "  $line"
	done < "$resource_log"
else
	echo "No created resources were logged."
fi

exit "$status"
