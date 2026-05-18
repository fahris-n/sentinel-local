# Project Sentinel

An API gateway built from scratch in Go that handles request routing, JWT authentication, role-based access control, and rate limiting for backend services. Deployed to a local Kubernetes cluster via minikube.

This is a personal learning project. The goal is to deeply understand API gateway internals, reverse proxy mechanics, middleware patterns, and cloud-native infrastructure rather than relying on off-the-shelf solutions like Kong or NGINX.

> **📖 For detailed package-by-package documentation, see the [Wiki](https://github.com/fahris-n/sentinel-local/wiki).**

---

## Architecture

```
Client Request
      │
      ▼
┌──────────────┐
│   Logging    │  ← Captures method, path, status code, request duration
├──────────────┤
│ Rate Limiter │  ← Token bucket algorithm (Lua + Redis), per-IP enforcement
├──────────────┤
│    Auth      │  ← JWT validation, role-based access control
├──────────────┤
│   Handler    │  ← Route matching → reverse proxy to backend
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Backends   │  ← Backend services (service-a, service-b, etc.)
└──────────────┘
```

Requests flow through a middleware chain where each layer can reject early. Rate-limited requests never touch auth, and unauthenticated requests never reach the backend.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Gateway | Go, `net/http`, `httputil.ReverseProxy` |
| Authentication | JWT (`golang-jwt/v5`), role-based access control |
| Rate Limiting | Token bucket algorithm, Lua scripting, Redis |
| Configuration | YAML-based route config with startup validation |
| Containerization | Docker, Docker Compose, multi-stage builds |
| Orchestration | Kubernetes (minikube), Deployments, Services |

---

## Features

**Request Routing** — YAML-configured routes map incoming paths to backend services via reverse proxies. Each route defines its own backend target, auth requirements, and rate limit parameters.

**JWT Authentication & RBAC** — Bearer tokens are validated per-request. Routes can require authentication and restrict access to specific roles (e.g., `user`, `premium`, `admin`). Public routes skip auth entirely.

**Token Bucket Rate Limiting** — Per-IP rate limiting using a token bucket algorithm implemented as an atomic Lua script executed inside Redis. Each route configures its own bucket size and refill rate.

**Middleware Chain** — Composable middleware pattern where middlewares wrap the handler in reverse order. Logging → Rate Limiting → Auth → Handler. Each layer can short-circuit the request.

**Config Validation** — Route configuration is validated at startup. Missing paths, empty backends, invalid token counts, and misconfigured auth/role combinations are caught before the gateway serves traffic.

**Kubernetes Deployment** — Full stack deployed to a local Kubernetes cluster via minikube. Services provide stable DNS-based discovery between the gateway, backends, and Redis.

---

## Project Structure

```
sentinel-local/
├── cmd/gateway/          # Application entry point
├── internal/
│   ├── auth/             # JWT validation, claims
│   ├── config/           # YAML loading, struct definitions, validation
│   ├── gateway/          # HTTP handler, route matching
│   ├── middleware/       # Logging, auth, rate limiting, chain
│   ├── proxy/            # Reverse proxy creation
│   ├── ratelimit/        # Redis connection, Lua script loader, limiter
│   └── routing/          # RouteEntry struct (runtime route representation)
├── configs/              # config.yaml
├── lua/                  # Token bucket Lua script
├── services/             # Backend service placeholders
├── k8s/                  # Kubernetes deployment & service manifests
├── docker-compose.yaml
└── Dockerfile
```

---

## Running Locally

### Docker Compose

```bash
docker compose up --build
```

### Kubernetes (minikube)

```bash
minikube start
minikube image load sentinel-local-gateway:latest
minikube image load sentinel-local-backend_a:latest
minikube image load sentinel-local-backend_b:latest
kubectl apply -f k8s/
minikube service gateway-service --url
```

### Testing

```bash
# Single request
curl -i http://<gateway-url>/api/register \
  -H "Authorization: Bearer <your-jwt>"

# Rate limit test
for i in $(seq 1 200); do
  curl -s -o /dev/null -w "%{http_code}\n" \
    http://<gateway-url>/api/register \
    -H "Authorization: Bearer <your-jwt>"
done
```

---

## What I Learned

### Why return errors instead of calling `log.Fatal` inside utility functions? 
Utility functions should return errors rather than throwing them so that the main application logic may handle the errors in any way deemed fit. 
For example, if the utility function for querying the Redis db fails an error is returned, allowing the main application logic to handle said error gracefully. 

In that example, an `Internal Server Error` is returned to the request since the gateway subscribes to the Fail-Closed pattern. This flexibility to handle errors in any way we choose is taken away when a program crashes due to errors being thrown instead of returned inside utility functions.

### How does Go's middleware chain pattern work? Why do the middlewares wrap in reverse order?

The middleware chain pattern allows for our `http.Handler` object to be wrapped in any number of middleware objects of our creation. This allows for the easy and clean handling of things such as logging, rate-limiting, and authentication. 

Middlewares are wrapped in reverse order so that the first middleware in the list is the outermost function. Think of the middlewares as being boxes, and the object being wrapped a toy. The last box in the list of boxes gets wrapped first so it's the closest to the toy, and the first box in the list gets wrapped last so its the outermost layer, the first reachable layer.

An example of this is shown in the `Chain.go` section of the [Middleware Docs](https://github.com/fahris-n/sentinel-local/wiki/Middleware#chaingo).

### What is the closure pattern in Go middleware, and why does `AuthMiddleware` have three nested layers?

Go supports anonymous functions (functions not bound to identifiers) and these anonymous functions can form [closures](https://go.dev/tour/moretypes/25). Closure functions are functions that reference variables from outside their function body.

[`AuthMiddleware`](https://github.com/fahris-n/sentinel-local/wiki/Middleware#auth-middlewarego) has three nested layers since all three layers are ran at different times and the information they capture becomes available at different time. Layer 1 runs once at startup and captures the route map. Layer 2 runs when the handler is wrapped in the `AuthMiddleware` and captures the next http handler. Layer 3 runs for every user request that hits the `AuthMiddleware` and has access to both the route map and next http handler captured above it. 

### What is the difference between `RouteConfig` and `RouteEntry`, and why are both needed?

_Your answer here._

### Why use a Lua script for rate limiting instead of handling the logic in Go with separate Redis commands?

_Your answer here._

### What does "fail open" vs "fail closed" mean, and which did you choose for rate limiting?

_Your answer here._

### Why rate limit before auth in the middleware chain?

_Your answer here._

### What problem does a Kubernetes Service solve that a Deployment alone doesn't?

_Your answer here._

### Why should Redis have only one replica in this architecture?

_Your answer here._

### What's the difference between how containers communicate in Docker Compose vs Kubernetes?

_Your answer here._

---

## Future Roadmap

- [ ] Kubernetes Secrets for JWT secret and Redis password
- [ ] Health checks (liveness/readiness probes)
- [ ] Structured JSON logging with request IDs
- [ ] Pattern-based route matching (`/api/users/*`)
- [ ] Request timeouts and backend resilience
- [ ] AWS deployment (EKS, Terraform, ECR)
