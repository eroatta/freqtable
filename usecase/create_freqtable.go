package usecase

import "github.com/eroatta/freqtable/entity"

type FrequencyTableUseCase interface {
	Create() (entity.FrequencyTable, error)
	Merge() (entity.FrequencyTable, error)
}

type asd struct {
}

/*
ft := FrequencyTable{
		ID: "",
	}

	values, err := fts.wcr.Extract(url)
	if err != nil {
		// TODO: log
		ft.Error = err
		return ft
	}
	ft.Values = values

	context := context.Background() // this is empty right now
	err = fts.ftr.Save(context, ft)
	if err != nil {
		// TODO: log
		ft.Error = err
		return ft
	}

	return ft
*/
