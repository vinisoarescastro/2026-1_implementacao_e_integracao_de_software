# Sistema Runner

> Trabalho Prático — Disciplina de Implementação e Integração (UFG, 2026-01)

---

## Visão Geral

O **Sistema Runner** é um conjunto de ferramentas de linha de comando (CLI) que permite executar aplicações Java sem que o usuário precise conhecer detalhes de configuração do ambiente Java. O sistema é composto por:

| Componente        | Tecnologia | Descrição                                                                    |
|-------------------|------------|------------------------------------------------------------------------------|
| `assinatura`      | Go 1.25    | CLI multiplataforma para criação e validação de assinaturas digitais         |
| `simulador`       | Go 1.25    | CLI multiplataforma para gerenciar o ciclo de vida do Simulador do HubSaúde  |
| `assinador.jar`   | Java 21    | Serviço de assinatura/validação (modo local e servidor HTTP)                 |

---

## Estrutura do Repositório

```
2026-1_implementacao_e_integracao_de_software/
│
├── cmd/
│   ├── assinatura/                  ← Binário principal do CLI de assinatura
│   │   ├── main.go
│   │   └── version_test.go
│   └── simulador/                   ← Binário do CLI do simulador (Sprint 4)
│       └── main.go
│
├── internal/
│   ├── cli/                         ← Parsing de comandos com Cobra
│   │   └── root.go
│   ├── invoker/                     ← Invocação do assinador.jar (local e HTTP)
│   │   ├── local.go
│   │   └── http.go
│   ├── jdk/                         ← Detecção e provisionamento automático do JDK
│   │   └── jdk.go
│   └── release/                     ← Download de artefatos (simulador.jar, JDK)
│       └── download.go
│
├── assinador/                       ← Projeto Java/Maven (assinador.jar)
│   ├── pom.xml
│   └── src/
│       ├── main/
│       │   └── java/com/kyriosdata/assinador/
│       │       ├── domain/
│       │       │   ├── SignRequest.java
│       │       │   ├── SignResponse.java
│       │       │   ├── ValidateRequest.java
│       │       │   └── ValidateResponse.java
│       │       ├── service/
│       │       │   ├── SignatureService.java       ← Interface principal
│       │       │   └── FakeSignatureService.java   ← Implementação simulada
│       │       └── Main.java
│       └── test/
│           └── java/com/kyriosdata/assinador/
│               └── service/
│                   └── FakeSignatureServiceTest.java
│
├── .github/
│   └── workflows/
│       ├── build.yml                ← CI: compila e testa em cada push
│       └── release.yml              ← CD: publica no GitHub Releases ao criar tag
│
├── orientacoes_trabalho_disciplina/ ← Documentação de orientação do professor
│   ├── README.md
│   ├── especificacao.md             ← Requisitos funcionais (US-01 a US-05)
│   ├── design.md                    ← Arquitetura C4 (contexto e contêineres)
│   ├── diagramas/
│   │   └── imagens/
│   │       ├── contexto.svg
│   │       └── conteineres.svg
│   ├── docs/
│   │   ├── plano-preliminar.md
│   │   ├── plano-revisitado.md
│   │   ├── plano-revisitado-v2.md   ← Plano oficial com 4 sprints
│   │   ├── planejamento.md
│   │   ├── sprint-1-tasks.md
│   │   └── implementacao_transcricao.md
│   └── projetos/
│       └── assinador-java/          ← Projeto Java base fornecido
│
├── docs/
│   └── manual-usuario.md            ← Manual de uso dos CLIs (entregável)
│
├── go.mod
├── go.sum
├── .gitignore
└── README.md                        ← Este arquivo
```

---

## Pré-requisitos

- [Go 1.25+](https://go.dev/dl/)
- [Java 21 (JDK)](https://adoptium.net/) - ou deixe o CLI provisionar automaticamente
- [Maven 3.9+](https://maven.apache.org/) - para compilar o assinador.jar
- [Git](https://git-scm.com/)
- [Cosign](https://docs.sigstore.dev/cosign/system_config/installation/) - para verificar artefatos (opcional)

---

## Como compilar

### CLI Go

```bash
# Compilar todos os binários
go build ./...

# Verificar problemas
go vet ./...

# Executar testes
go test ./...
```

### assinador.jar (Java)

```bash
cd assinador/
mvn clean package
```

### Build multiplataforma (via CI ou local)

```bash
# Linux
GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.version=v1.0.0" -o dist/assinatura-v1.0.0-linux-amd64   ./cmd/assinatura

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=v1.0.0" -o dist/assinatura-v1.0.0-windows-amd64.exe ./cmd/assinatura

# macOS
GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.version=v1.0.0" -o dist/assinatura-v1.0.0-darwin-amd64  ./cmd/assinatura
```

---

## Uso

### CLI `assinatura`

```bash
# Exibir versão
assinatura version

# Criar assinatura (modo local — invoca java -jar diretamente)
assinatura sign --content "documento.pdf" --token "meu-token" --local

# Criar assinatura (modo servidor — usa instância HTTP em execução)
assinatura sign --content "documento.pdf" --token "meu-token"

# Validar assinatura
assinatura validate --content "documento.pdf" --signature "base64..."

# Iniciar o assinador.jar como servidor HTTP
assinatura start [--port 8080]

# Encerrar o servidor
assinatura stop [--port 8080]

# Encerrar após inatividade
assinatura start --timeout 30
```

### CLI `simulador`

```bash
# Iniciar o Simulador do HubSaúde
simulador start [--source <url-alternativa>]

# Parar o simulador
simulador stop

# Verificar status
simulador status
```

---

## Roadmap — O que precisa ser feito

O projeto é desenvolvido em **4 sprints de 1 semana** com estratégia iterativa e incremental.

### ✅ Sprint 1 — Fundação & CI/CD *(concluída)*

- [x] Módulo Go inicializado (`go mod init github.com/kyriosdata/runner`)
- [x] Estrutura de pacotes criada conforme DT-06
- [x] Comando `assinatura version` funcionando com Cobra
- [x] Stub do binário `simulador` criado
- [x] GitHub Actions: build multiplataforma (Windows, Linux, macOS) a cada push
- [x] Workflow de release com SemVer ao criar tag `v*`
- [x] Checksums SHA-256 gerados para cada binário
- [x] Artefatos assinados com Cosign/Sigstore (keyless, via OIDC)

---

### 🔄 Sprint 2 — Assinatura Simulada, modo local *(em andamento)*

**Objetivo:** fluxo ponta-a-ponta funcional — usuário executa `assinatura sign` e obtém resultado sem configurar Java.

- [x] Projeto Java base criado com `pom.xml` e estrutura Maven
- [x] Interface `SignatureService` definida com métodos `sign` e `validate`
- [x] `FakeSignatureService` retorna assinatura simulada para parâmetros válidos
- [x] Testes unitários da `FakeSignatureService` passando
- [ ] **US-02.2** — Validação completa dos parâmetros de criação de assinatura no `assinador.jar`
  - Verificar presença e formato de todos os campos obrigatórios
  - Retornar mensagem de erro indicando qual parâmetro está inválido e o motivo
- [ ] **US-02.3** — Validação de parâmetros e simulação do fluxo de *validação* de assinatura
- [ ] **US-01.2** — Comandos `sign` e `validate` no CLI Go com Cobra
  - Flags mapeados para os parâmetros do jar
  - `--help` documentado
- [ ] **US-01.3** — Invocação do `assinador.jar` via `java -jar` no pacote `internal/invoker`
  - Usar `os/exec` para invocar o processo
  - Capturar stdout/stderr e repassar ao usuário
  - Tratar erros: JDK ausente, jar não encontrado, parâmetros inválidos
- [ ] **US-01.4** — Exibição legível dos resultados no terminal
  - Sucesso: exibe campos da assinatura de forma estruturada
  - Erro: indica claramente o problema e como corrigir
- [ ] **US-04.1** — Detecção e provisionamento automático do JDK 21
  - Verificar se `java` está disponível no `PATH` ou em `~/.hubsaude/jdk/`
  - Se ausente, baixar Temurin/Zulu para a plataforma correta
  - Armazenar em `~/.hubsaude/jdk/` e reutilizar nas próximas execuções

---

### ⏳ Sprint 3 — Modo Servidor HTTP & PKCS#11

**Objetivo:** o `assinador.jar` roda como servidor HTTP; o CLI gerencia seu ciclo de vida e se comunica via HTTP.

- [ ] **US-02.4** — Endpoints HTTP no `assinador.jar`
  - `POST /sign` e `POST /validate`
  - Reutilizar a lógica da `FakeSignatureService`
  - Respostas HTTP com estrutura consistente (sucesso e erro)
  - Testes de integração dos endpoints
- [ ] **US-02.5** — Integração com dispositivo criptográfico via PKCS#11
  - Usar `SunPKCS11` do Java
  - Simulação com [SoftHSM2](https://github.com/softhsm/SoftHSMv2)
- [ ] **US-01.5** — CLI inicia o `assinador.jar` no modo servidor
  - Detectar porta disponível se a padrão estiver ocupada
  - Registrar PID e porta em `~/.hubsaude/`
- [ ] **US-01.6** — CLI detecta instância já em execução e a reutiliza
- [ ] **US-01.7** — CLI envia requisições HTTP ao servidor quando em modo servidor
- [ ] **US-01.8** — Comando `assinatura stop` encerra o servidor
- [ ] **US-01.9** — Suporte a `--timeout <minutos>` para encerramento por inatividade

---

### ⏳ Sprint 4 — CLI Simulador & Entrega Final

**Objetivo:** sistema completo com gestão do Simulador do HubSaúde.

- [ ] **US-03.1** — `simulador start` inicia o `simulador.jar`
  - Verificar portas disponíveis antes de iniciar
  - Baixar `simulador.jar` automaticamente se não estiver localmente disponível
- [ ] **US-03.2** — `simulador stop` e `simulador status`
  - Registrar PID e porta em `~/.hubsaude/`
  - Encerramento limpo com tratamento de erros
- [ ] **US-03.3** — CLI `simulador` completo com pipeline CI/CD gerando binários multiplataforma
- [ ] **US-03.4** — Download automático do `simulador.jar` do GitHub Releases
  - Opção `--source <url>` para URL alternativa
  - Cache local: não re-baixar se versão já disponível
  - Verificação de integridade (checksum SHA-256)

---

## Entregáveis obrigatórios

Conforme `orientacoes_trabalho_disciplina/especificacao.md` (seção 7):

| # | Entregável | Status |
|---|---|---|
| 1 | Código-fonte do CLI `assinatura` (Go, multiplataforma) | 🔄 Em andamento |
| 2 | Código-fonte do `assinador.jar` (Java 21, Maven) | 🔄 Em andamento |
| 3 | Testes (unitários, integração, aceitação, cenários de erro) | 🔄 Em andamento |
| 4 | Documentação (manual de usuário, guia técnico, exemplos, instalação) | ⏳ Pendente |
| 5 | Especificação com diagramas C4 | ✅ Fornecida pelo professor |
| 6 | Binários pré-compilados para Win/Linux/macOS via GitHub Releases | ✅ Automatizado |
| 7 | Código-fonte do CLI `simulador` (Go, multiplataforma) | ⏳ Sprint 4 |

---

## Arquitetura

O sistema segue o **Modelo C4**. Os diagramas estão em `orientacoes_trabalho_disciplina/diagramas/imagens/`.

### Diagrama de Contêineres (resumido)

```
Usuário
  │
  ├──(CLI)──► assinatura         (Go 1.25)
  │               │
  │               ├──(java -jar / HTTP)──► assinador.jar   (Java 21)
  │               │                             │
  │               │                        (PKCS#11)
  │               │                             │
  │               │                    Dispositivo Criptográfico
  │               │                    (token USB / SoftHSM2)
  │
  └──(CLI)──► simulador          (Go 1.25)
                  │
               (HTTP)
                  │
          Simulador do HubSaúde  (sistema externo)
```

### Dados locais (`~/.hubsaude/`)

```
~/.hubsaude/
├── jdk/            ← JDK 21 provisionado automaticamente
├── simulador.jar   ← baixado dinamicamente
└── state.json      ← PID, porta e estado dos processos em execução
```

---

## CI/CD

| Workflow | Gatilho | O que faz |
|---|---|---|
| `build.yml` | Push na `main` | Compila, executa `go vet`, roda testes |
| `release.yml` | Criação de tag `v*` | Compila para 3 plataformas, gera checksums SHA-256, assina com Cosign, publica no GitHub Releases |

### Verificar autenticidade de um artefato

```bash
cosign verify-blob \
  --certificate assinatura-v1.0.0-linux-amd64.pem \
  --signature  assinatura-v1.0.0-linux-amd64.sig \
  --certificate-identity "https://github.com/<seu-usuario>/runner/.github/workflows/ci.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  assinatura-v1.0.0-linux-amd64
```

---

## Decisões técnicas

| # | Decisão | Valor |
|---|---|---|
| DT-01 | Módulo Go | `github.com/kyriosdata/runner` |
| DT-02 | Branch principal | `main` |
| DT-03 | Plataformas-alvo | `windows/amd64`, `linux/amd64`, `darwin/amd64` |
| DT-04 | Convenção de artefatos | `assinatura-<versão>-<os>-<arch>` |
| DT-05 | Checksums | SHA-256 por binário + arquivo `checksums.txt` no release |
| DT-06 | Layout de pacotes | `cmd/` para binários, `internal/` para pacotes compartilhados |

---

## Referências

- [Especificação do projeto](orientacoes_trabalho_disciplina/especificacao.md)
- [Design e arquitetura C4](orientacoes_trabalho_disciplina/design.md)
- [Plano de implementação (Sprints)](orientacoes_trabalho_disciplina/docs/plano-revisitado-v2.md)
- [Casos de uso — Criar Assinatura (FHIR HubSaúde)](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-criar-assinatura.html)
- [Casos de uso — Validar Assinatura (FHIR HubSaúde)](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-validar-assinatura.html)
- [Cobra CLI (Go)](https://cobra.dev/)
- [SoftHSM2 — simulador PKCS#11](https://github.com/softhsm/SoftHSMv2)
- [Sigstore / Cosign](https://docs.sigstore.dev/)
- [Modelo C4](https://c4model.com/)