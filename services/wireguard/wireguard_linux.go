package wireguard

func (w *WireGuard) RealInterface() (string, error) {
	return w.cfg.Name, nil
}
