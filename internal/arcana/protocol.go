package arcana

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

const protocolID = "arcana.query.v1"

type request struct {
	ID               string   `json:"id"`
	Op               string   `json:"op"`
	PathPrefix       string   `json:"path_prefix,omitempty"`
	Path             string   `json:"path,omitempty"`
	Name             string   `json:"name,omitempty"`
	Kind             string   `json:"kind,omitempty"`
	Relation         string   `json:"relation,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Offset           int      `json:"offset,omitempty"`
	NodeID           uint32   `json:"node_id"`
	Direction        string   `json:"direction,omitempty"`
	Relations        []string `json:"relations,omitempty"`
	MinCommunitySize int      `json:"min_community_size,omitempty"`
}

type response struct {
	Protocol string          `json:"protocol"`
	ID       string          `json:"id"`
	OK       bool            `json:"ok"`
	Result   json.RawMessage `json:"result"`
	Error    *responseError  `json:"error,omitempty"`
}

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type protocolRun func(context.Context, string, string, []request) (map[string]response, error)

func runProtocol(
	ctx context.Context,
	command string,
	snapshot string,
	requests []request,
) (map[string]response, error) {
	var input bytes.Buffer
	encoder := json.NewEncoder(&input)
	for _, current := range requests {
		if err := encoder.Encode(current); err != nil {
			return nil, fmt.Errorf("encode Arcana request: %w", err)
		}
	}

	var stdout, stderr bytes.Buffer
	process := exec.CommandContext(ctx, command, "protocol", "--snapshot", snapshot)
	process.Stdin = &input
	process.Stdout = &stdout
	process.Stderr = &stderr
	if err := process.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message != "" {
			return nil, fmt.Errorf("run Arcana protocol: %w: %s", err, message)
		}
		return nil, fmt.Errorf("run Arcana protocol: %w", err)
	}

	responses := make(map[string]response, len(requests))
	scanner := bufio.NewScanner(&stdout)
	scanner.Buffer(make([]byte, 64*1024), 64*1024*1024)
	for scanner.Scan() {
		var decoded response
		if err := json.Unmarshal(scanner.Bytes(), &decoded); err != nil {
			return nil, fmt.Errorf("decode Arcana response: %w", err)
		}
		if decoded.Protocol != protocolID {
			return nil, fmt.Errorf("unexpected Arcana protocol %q", decoded.Protocol)
		}
		responses[decoded.ID] = decoded
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read Arcana responses: %w", err)
	}
	if len(responses) != len(requests) {
		return nil, fmt.Errorf("Arcana returned %d responses for %d requests", len(responses), len(requests))
	}
	for _, current := range requests {
		decoded, exists := responses[current.ID]
		if !exists {
			return nil, fmt.Errorf("Arcana did not return response %q", current.ID)
		}
		if decoded.OK {
			continue
		}
		if decoded.Error == nil {
			return nil, fmt.Errorf("Arcana operation %s failed without an error payload", current.Op)
		}
		return nil, fmt.Errorf(
			"Arcana operation %s failed (%s): %s",
			current.Op,
			decoded.Error.Code,
			decoded.Error.Message,
		)
	}
	return responses, nil
}
