package log

type RpcLogger int

func (r *RpcLogger) ListAll(_, reply *[]*LoggerEntry) error {
	*reply = ListAll()
	return nil
}

func (r *RpcLogger) SetLevel(args NamedLogger, reply *bool) error {
	if level, err := parseLevel(args.Level); err != nil {
		return err
	} else {
		*reply = SetLevel(args.Name, level)
		return nil
	}
}

func (r *RpcLogger) SetLevelToTree(args NamedLogger, reply *int) error {
	if level, err := parseLevel(args.Level); err != nil {
		return err
	} else {
		*reply = SetLevelToTree(args.Name, level)
		return nil
	}
}
