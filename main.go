package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dadosjusbr/coletores/status"
	"github.com/dadosjusbr/proto/pipeline"
	"github.com/frictionlessdata/datapackage-go/datapackage"

	"google.golang.org/protobuf/encoding/prototext"

	frictionless "github.com/frictionlessdata/tableschema-go/csv"
)

var resources = []string{"coleta", "contra_cheque", "remuneracao"}

func main() {
	// Reading input.
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		err := status.NewError(status.InvalidInput, fmt.Errorf("error reading execution result from stdin: %q", err))
		status.ExitFromError(err)
	}
	var er pipeline.ResultadoExecucao
	if err := prototext.Unmarshal(in, &er); err != nil {
		err := status.NewError(status.InvalidInput, fmt.Errorf("error unmarshalling execution result: %q", err))
		status.ExitFromError(err)
	}

	// Loading and validating package.
	pkg, err := datapackage.Load(er.Pr.Pacote)
	if err != nil {
		err = status.NewError(status.DataUnavailable, fmt.Errorf("error loading datapackage (%s):%q", er.Pr.Pacote, err))
		status.ExitFromError(err)
	}

	for _, v := range resources {
		sch, err := pkg.GetResource(v).GetSchema()
		if err != nil {
			err = status.NewError(status.DataUnavailable, fmt.Errorf("error getting schema from data package resource (%s | %s):%q", er.Pr.Pacote, v, err))
			status.ExitFromError(err)
		}
		if err := sch.Validate(); err != nil {
			err = status.NewError(status.InvalidInput, fmt.Errorf("error validating schema (%s):%q", er.Pr.Pacote, err))
			status.ExitFromError(err)
		}

		switch v {
		case "coleta":
			if err := pkg.GetResource(v).Cast(&[]Coleta_CSV{}, frictionless.LoadHeaders()); err != nil {
				err = status.NewError(status.InvalidInput, fmt.Errorf("error validating datapackage (%s):%q", er.Pr.Pacote, err))
				status.ExitFromError(err)
			}
		case "remuneracao":
			if err := pkg.GetResource(v).Cast(&[]Remuneracao_CSV{}, frictionless.LoadHeaders()); err != nil {
				err = status.NewError(status.InvalidInput, fmt.Errorf("error validating datapackage (%s):%q", er.Pr.Pacote, err))
				status.ExitFromError(err)
			}
		case "metadados":
			if err := pkg.GetResource(v).Cast(&[]Metadados_CSV{}, frictionless.LoadHeaders()); err != nil {
				err = status.NewError(status.InvalidInput, fmt.Errorf("error validating datapackage (%s):%q", er.Pr.Pacote, err))
				status.ExitFromError(err)
			}
		case "contra_cheque":
			if err := pkg.GetResource(v).Cast(&[]ContraCheque_CSV{}, frictionless.LoadHeaders()); err != nil {
				err = status.NewError(status.InvalidInput, fmt.Errorf("error validating datapackage (%s):%q", er.Pr.Pacote, err))
				status.ExitFromError(err)
			}
		default:
			if err := pkg.GetResource(v).Cast(&[]ContraCheque_CSV{}, frictionless.LoadHeaders()); err != nil {
				err = status.NewError(status.InvalidInput, fmt.Errorf("error validating datapackage (%s):%q", er.Pr.Pacote, err))
				status.ExitFromError(err)
			}
		}

	}

	// Printing output.
	b, err := prototext.Marshal(&er)
	if err != nil {
		err = status.NewError(status.Unknown, fmt.Errorf("error marshalling output:%q", err))
		status.ExitFromError(err)
	}
	fmt.Printf("%s", b)
}
