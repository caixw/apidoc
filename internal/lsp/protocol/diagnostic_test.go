// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestPublishDiagnosticParams_AppendDiagnostic(t *testing.T) {
	a := assert.New(t)
	p := NewPublishDiagnosticsParams(core.URI("test.go"))
	a.NotNil(p).
		Equal(p.URI, core.URI("test.go")).
		Empty(p.Diagnostics)

	p.AppendDiagnostic(core.NewError(locale.ErrInvalidUTF8Character), core.Erro)
	a.Equal(1, len(p.Diagnostics))
	p.AppendDiagnostic(core.NewError(locale.ErrInvalidUTF8Character), core.Warn)
	a.Equal(2, len(p.Diagnostics))
	p.AppendDiagnostic(core.NewError(locale.ErrInvalidUTF8Character), core.Info)
	a.Equal(3, len(p.Diagnostics))
	// 忽略 core.Succ
	p.AppendDiagnostic(core.NewError(locale.ErrInvalidUTF8Character), core.Succ)
	a.Equal(3, len(p.Diagnostics))

	a.Panic(func() {
		p.AppendDiagnostic(core.NewError(locale.ErrInvalidUTF8Character), 100)
	})
}

func TestBuildDiagnostic(t *testing.T) {
	a := assert.New(t)

	err := core.NewError(locale.ErrInvalidUTF8Character).WithLocation(core.Location{
		Range: core.Range{Start: core.Position{Line: 1}},
	})
	d := buildDiagnostic(err, DiagnosticSeverityWarning)
	a.Equal(d.Severity, DiagnosticSeverityWarning)
	a.Empty(d.Tags)
	a.Equal(d.Range.Start.Line, 1)

	err = err.AddTypes(core.ErrorTypeDeprecated, core.ErrorTypeUnused)
	d = buildDiagnostic(err, DiagnosticSeverityError)
	a.Equal(d.Severity, DiagnosticSeverityError)
	a.Equal(d.Tags, []DiagnosticTag{DiagnosticTagDeprecated, DiagnosticTagUnnecessary})
	a.Equal(d.Range.Start.Line, 1)
}
