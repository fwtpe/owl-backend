package log

type RpcLogger int

func (r *RpcLogger) ListAll(_, reply []*LoggerEntry) error {
	reply = ListAll()
	return nil
}

func (r *RpcLogger) SetLevel(args NamedLogger, reply *bool) error {
	if level, err := parseLevel(args.Level); err != nil {
		return err
	} else {
		ret := SetLevel(args.Name, level)
		reply = &ret
		return nil
	}
}

func (r *RpcLogger) SetLevelToTree(args NamedLogger, reply *int) error {
	if level, err := parseLevel(args.Level); err != nil {
		return err
	} else {
		ret := SetLevelToTree(args.Name, level)
		reply = &ret
		return nil
	}
}
