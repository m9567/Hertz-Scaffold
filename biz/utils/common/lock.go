package common

import "github.com/duke-git/lancet/v2/concurrency"

var Locker *concurrency.TryKeyedLocker[string]
