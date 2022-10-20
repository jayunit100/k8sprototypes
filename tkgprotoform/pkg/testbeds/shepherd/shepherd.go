package shepherd

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.eng.vmware.com/shepherd/shepherd/cli/api"
	"gitlab.eng.vmware.com/shepherd/shepherd/cli/cmd"
	"gitlab.eng.vmware.com/shepherd/shepherd/cli/common"
	"gitlab.eng.vmware.com/shepherd/shepherd/cli/lock"
	"log"
	"os"
	"time"
)

const lockCreateDesc = ``

func IdempotentBootstrap() {
	if _, err := os.Stat("~/corgo"); err != nil {
		os.Mkdir("~/corgo", 0755)
	}
}

func defaultTarget() cmd.Target {
	return cmd.ViperTarget{}
}

// Shepherd implements TestBed interface
type Shepherd struct {
	Lock   string
	Target cmd.Target
}

// CreateTestbed returns a Shepherd impl  of the Testbed.
// See https://gitlab.eng.vmware.com/shepherd/shepherd/-/blob/main/cli/cmd/lock_create.go#L22
// for the base impl for this...
func (s *Shepherd) Create(target cmd.Target) (lockid string, err error) {

	if target == nil {
		s.Target = defaultTarget()
	}
	params := lock.CreateLockParams{}
	params.Annotations = make(common.Annotations)
	params.Recipe.Data = make(map[string]string)
	params.Recipe.Meta = make(map[string]string)

	IdempotentBootstrap()
	myLock := lock.LockCreationParams{
		QuotaTimeout: 5 * time.Minute,
		Block:        false,
		OutputPath:   "~/corgo/corgo.json",
		OutputStream: os.Stdout,
		CreateLock: func(ctx context.Context) (lockid string, err error) {
			log.Printf("Creating a lock:")
			lockId, _, err := lock.CreateLock(ctx, target.Client(), target.Namespace(), params)
			// set the lock of the struct
			s.Lock = lockId
			return s.Lock, errors.WithMessage(api.ExpandError(err), "failed to create lock")
		},
	}
	ctx := context.Background()
	return myLock.CreateLock(ctx)

}
