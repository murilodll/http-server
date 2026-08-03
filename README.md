# http-server [GO] + Nginx + Prometheus + Grafana 

Servidor HTTP em Go, containerizado, com stack de observabilidade (Prometheus + Grafana) e reverse proxy via Nginx. Inclui também um playbook Ansible para provisionamento/deploy automatizado.

## Sumário

- [Visão geral](#visão-geral)
- [Arquitetura](#arquitetura)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Endpoints](#endpoints)
- [Como rodar](#como-rodar)
- [Observabilidade](#observabilidade)
- [Deploy com Ansible](#deploy-com-ansible)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Requisitos](#requisitos)

## Visão geral

O `http-server-projeto-korp` é uma aplicação Go que expõe um endpoint simples retornando nome do projeto e horário atual em JSON, além de endpoints de health check e métricas Prometheus. Toda a stack roda em containers Docker orquestrados via `docker-compose`, com Nginx como ponto de entrada único (porta 80), Prometheus coletando métricas do servidor e Grafana exibindo dashboards já provisionados.

## Arquitetura

```
Cliente
  │
  ▼
Nginx (porta 80)
  ├── /projeto-korp  → http-server-projeto-korp:8080
  ├── /health        → http-server-projeto-korp:8080
  ├── /metrics       → http-server-projeto-korp:8080
  ├── /grafana/      → grafana:3000
  └── /prometheus/   → prometheus:9090

Prometheus ──(scrape a cada 5s)──> http-server-projeto-korp:8080/metrics
Grafana ──(datasource)──> Prometheus
```

Todos os serviços compartilham a rede Docker `korp-network`.

## Estrutura do projeto

```
.
├── cmd/
│   └── server/
│       └── main.go              # entrypoint: rotas, middleware de métricas, graceful shutdown
├── internal/
│   ├── handler/
│   │   └── handler.go           # handlers HTTP (/projeto-korp e /health)
│   └── model/
│       └── model.go             # struct de resposta (Response)
├── nginx/
│   └── http-server-projeto-korp.conf   # config do reverse proxy
├── prometheus/
│   └── prometheus.yml           # config de scrape do Prometheus
├── grafana/
│   └── provisioning/
│       ├── datasources/         # datasource do Prometheus pré-configurado
│       └── dashboards/          # dashboard "korp" pré-configurado
├── ansible/
│   ├── playbook.yml             # playbook de provisionamento/deploy
│   └── ansible.cfg
├── Dockerfile                   # build multi-stage (Go → Alpine)
├── docker-compose.yml           # orquestração dos serviços
├── go.mod / go.sum
└── README.md
```

## Endpoints

| Método | Rota            | Descrição                                             |
|--------|-----------------|--------------------------------------------------------|
| GET    | `/projeto-korp` | Retorna JSON `{ "nome": "Projeto Korp", "horario": "<RFC3339 UTC>" }` |
| GET    | `/health`       | Health check simples, retorna `OK` (texto puro)        |
| GET    | `/metrics`      | Métricas no formato Prometheus (inclui `http_requests_total`) |

Exemplo de resposta de `/projeto-korp`:

```json
{
  "nome": "Projeto Korp",
  "horario": "2026-08-03T12:00:00Z"
}
```

## Como rodar

### Com Docker Compose (recomendado)

```bash
# defina a variável usada pelo Grafana para montar a Root URL
echo "SERVER_IP=localhost" > .env

docker compose up -d --build
```

Após subir, os serviços ficam acessíveis via Nginx:

- Aplicação: `http://localhost/projeto-korp`
- Health check: `http://localhost/health`
- Métricas: `http://localhost/metrics`
- Grafana: `http://localhost/grafana/`
- Prometheus: `http://localhost/prometheus/`

Para derrubar o ambiente:

```bash
docker compose down
```

### Localmente, sem Docker

```bash
go run ./cmd/server
```

O servidor sobe diretamente na porta `8080` (sem Nginx/Prometheus/Grafana).

## Observabilidade

- **Métricas**: o servidor expõe `/metrics` via `promhttp`, incluindo o contador customizado `http_requests_total`, rotulado por rota (`path`), incrementado a cada requisição.
- **Prometheus**: configurado para fazer scrape de `http-server-projeto-korp:8080` a cada 5 segundos (`prometheus/prometheus.yml`).
- **Grafana**: já vem com o datasource do Prometheus e um dashboard (`dashboard-korp.json`) provisionados automaticamente ao subir o container — não é necessário configurar nada manualmente.

## Deploy com Ansible

O diretório `ansible/` contém um playbook (`playbook.yml`) para automatizar o provisionamento do ambiente em um host remoto (ex: instalação de dependências, cópia dos arquivos do projeto, subida via Docker Compose).

```bash
cd ansible
ansible-playbook -i <seu_inventario> playbook.yml
```

> Ajuste o inventário e as variáveis conforme o ambiente de destino.

## Variáveis de ambiente

| Variável    | Usado por | Descrição                                              |
|-------------|-----------|----------------------------------------------------------|
| `SERVER_IP` | Grafana   | Define a `GF_SERVER_ROOT_URL` (`http://${SERVER_IP}/grafana/`) |

## Requisitos

- Docker e Docker Compose
- Go 1.25+ (apenas para rodar/buildar localmente fora de containers)
- Ansible (opcional, apenas para deploy automatizado)
