# ADR-001 — Modo servidor HTTP como padrão de invocação

## Contexto

O CLI `assinatura` pode invocar o `assinador.jar` de duas formas:

- **Modo local**: executa `java -jar assinador.jar` a cada chamada. Simples, mas paga o custo de inicialização da JVM (~200–500ms) em toda invocação.
- **Modo servidor**: o `assinador.jar` fica em execução contínua como servidor HTTP. O CLI envia requisições HTTP, eliminando o overhead de cold start nas chamadas subsequentes.

A especificação (US-01) define explicitamente que o modo servidor deve ser o padrão e o modo local deve ser ativado explicitamente.

## Decisão

O modo servidor HTTP é o modo padrão. O modo local é ativado pela flag `--local`.

```
assinatura sign --content "doc.pdf"          # usa HTTP (padrão)
assinatura sign --content "doc.pdf" --local  # invoca java -jar diretamente
```

A porta padrão do servidor é **8080** (ver ADR-003).

## Consequências

**Fica mais fácil:**
- Múltiplas operações em sequência têm latência significativamente menor.
- O CLI tem uma única responsabilidade de transporte por vez — não mistura lógica de processo e HTTP.

**Fica mais difícil:**
- O usuário precisa iniciar o servidor antes de usar o CLI em modo padrão (via `assinatura start`).
- Cenários de uso esporádico exigem o `--local` explícito enquanto o servidor não está em execução.
- O CLI precisa de lógica para detectar se o servidor está ativo e fornecer mensagem de erro orientativa.
