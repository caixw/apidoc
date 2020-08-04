// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestBuildDiagnostic(t *testing.T) {
	a := assert.New(t)

	err := core.NewError(locale.ErrRequired).WithLocation(core.Location{
		Range: core.Range{Start: core.Position{Line: 1}},
	})
	d := BuildDiagnostic(err, DiagnosticSeverityWarning)
	a.Equal(d.Severity, DiagnosticSeverityWarning)
	a.Empty(d.Tags)
	a.Equal(d.Range.Start.Line, 1)

	err.AddTypes(core.ErrorTypeDeprecated, core.ErrorTypeUnused)
	d = BuildDiagnostic(err, DiagnosticSeverityError)
	a.Equal(d.Severity, DiagnosticSeverityError)
	a.Equal(d.Tags, []DiagnosticTag{DiagnosticTagDeprecated, DiagnosticTagUnnecessary})
	a.Equal(d.Range.Start.Line, 1)
}
