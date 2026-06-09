# ADR-003 — Porta padrão 8080 para o assinador.jar

## Contexto

O `assinador.jar` em modo servidor precisa escutar em uma porta TCP. A porta deve:

- Não conflitar com serviços comuns em ambientes de desenvolvimento (80, 443, 3000, 5000, 8443).
- Não exigir privilégios de root (portas abaixo de 1024 exigem no Linux/macOS).
- Ser configurável via `--port` para ambientes onde a porta padrão já está ocupada.

A porta 8443 já está reservada para o Simulador do HubSaúde (especificação US-03).

## Decisão

A porta padrão do `assinador.jar` em modo servidor é **8080**.

O CLI respeita a flag `--port <número>` para sobrescrever o padrão tanto no `start` quanto no `stop`.

## Consequências

**Fica mais fácil:**
- 8080 é a porta de desenvolvimento HTTP mais reconhecida; desenvolvedores a identificam imediatamente.
- Não conflita com a porta 8443 reservada ao Simulador, permitindo ambos rodarem simultaneamente.
- Acima de 1024 — sem necessidade de privilégios elevados em nenhum SO suportado.

**Fica mais difícil:**
- Ambientes com servidor de aplicação local (Tomcat, WildFly) frequentemente usam 8080, exigindo uso de `--port`.
- O CLI precisa detectar porta ocupada e fornecer mensagem de erro clara orientando o uso de `--port`.
