package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"gear-contact/internal/gears"
	"gear-contact/internal/mesh"
	"gear-contact/internal/server"
)

const usage = `gear-contact: involute spur gear mesh calculator.

Reads a JSON gear pair and prints the engagement geometry: pitch / base /
addendum diameters, centre distance, line of action, contact ratio, tip
clearance and undercut status.

usage:
  gear-contact mesh <gear-spec.json>
  gear-contact -http [addr]
  gear-contact help

gear-spec.json fields:
  module                 module m in mm (strictly positive)
  pressure_angle_deg     standard pressure angle in degrees (14.5..25, < 90)
  addendum_coefficient   ha* (default 1)
  clearance_coefficient  c* (default 0.25)
  gear1                  {teeth: z1, shift: x1 (default 0)}
  gear2                  {teeth: z2, shift: x2 (default 0)}

Illegal inputs (module <= 0, teeth <= 0, pressure angle out of range, an
undercut gear without sufficient profile shift, a line of action that never
meets the addendum circles) are reported on stderr and exit non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "mesh":
		if err := runMesh(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "gear-contact: %v\n", err)
			os.Exit(1)
		}
	case "-http", "--http":
		if err := runHTTP(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "gear-contact: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "gear-contact: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runMesh(args []string) error {
	if len(args) != 1 {
		return errors.New("mesh needs exactly one gear-spec JSON file")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read spec: %w", err)
	}
	pair, err := gears.FromJSON(data)
	if err != nil {
		return fmt.Errorf("spec: %w", err)
	}
	m, err := mesh.Compute(pair)
	if err != nil {
		return fmt.Errorf("mesh: %w", err)
	}
	fmt.Print(m.RenderText())
	return nil
}

func runHTTP(args []string) error {
	addr := ":8080"
	if len(args) > 1 {
		return errors.New("-http takes at most one listen address")
	}
	if len(args) == 1 {
		addr = args[0]
	}
	fmt.Printf("gear-contact HTTP on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, server.New("web", "example/spur-20.json"))
}
