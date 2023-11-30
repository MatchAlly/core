package tournament

type TournamentFormat string

const (
	FormatSingleElimination TournamentFormat = "single_elimination"
	FormatDoubleElimination TournamentFormat = "double_elimination"
	FormatRoundRobin        TournamentFormat = "round_robin"
	FormatSwiss             TournamentFormat = "swiss"
)

func (s *ServiceImpl) createSingleEliminationTournament(teams [][]uint) (*Tournament, error) {
	panic("not implemented")
}

func (s *ServiceImpl) createDoubleEliminationTournament(teams [][]uint) (*Tournament, error) {
	panic("not implemented")
}

func (s *ServiceImpl) createRoundRobinTournament(teams [][]uint) (*Tournament, error) {
	panic("not implemented")
}

func (s *ServiceImpl) createSwissTournament(teams [][]uint) (*Tournament, error) {
	panic("not implemented")
}
