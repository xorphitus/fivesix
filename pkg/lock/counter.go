package lock

import (
	bpf "github.com/iovisor/gobpf/bcc"
)

const source = `
#include <uapi/linux/ptrace.h>

struct key_t {
    char c[80];
};
BPF_HASH(counts, struct key_t);

int increment(struct pt_regs *ctx) {
  	struct key_t key = {};
	u64 zero = 0, *val;

  	val = counts.lookup_or_init(&key, &zero);
	(*val)++;

	return 0;
};

int decrement(struct pt_regs *ctx) {
  	struct key_t key = {};
	u64 zero = 0, *val;

  	val = counts.lookup_or_init(&key, &zero);
	(*val)--;

	return 0;
};
`

type probe struct {
	module  *bpf.Module
	binPath string
	pid     int
}

func Run(binPath string, pid int) (table *bpf.Table, done func(), err error) {
	m := bpf.NewModule(source, []string{})

	p := &probe{
		module:  m,
		binPath: binPath,
		pid:     pid,
	}

	err = p.attach("internal/poll.(*fdMutex).rwlock", "increment")
	if err != nil {
		return nil, nil, err
	}

	err = p.attach("internal/poll.(*fdMutex).rwunlock", "decrement")
	if err != nil {
		return nil, nil, err
	}

	table = bpf.NewTable(m.TableId("counts"), m)

	return table, m.Close, nil
}

func (p *probe) attach(fnSym, probeFn string) error {
	fd, err := p.module.LoadUprobe(probeFn)
	if err != nil {
		return err
	}

	err = p.module.AttachUprobe(p.binPath, fnSym, fd, p.pid)
	if err != nil {
		return err
	}

	return nil
}
