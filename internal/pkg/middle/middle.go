package middle

import "github.com/mozgunovdm/example/internal/pkg/employe"

// Middle describes a service middleware.
type Middleware func(service employe.Service) employe.Service
