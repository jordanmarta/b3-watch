package decision

// PrecoAtualAbaixoDoPrecoMedio
// Regra básica: preço atual < PM contábil
func PrecoAtualAbaixoDoPrecoMedio(precoAtual, precoMedio float64) bool {
	return precoAtual < precoMedio
}

// PrecoAtualAbaixoDoPercentual
// Ex: percentual = 5 -> alerta quando precoAtual < precoMedio * 0.95
func PrecoAtualAbaixoDoPercentual(precoAtual, precoMedio, percentual float64) bool {
	if percentual <= 0 {
		return false
	}
	limite := precoMedio * (1 - percentual/100)
	return precoAtual < limite
}
