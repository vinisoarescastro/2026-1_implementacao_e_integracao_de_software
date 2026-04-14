package com.kyriosdata.assinador;

import com.kyriosdata.assinador.domain.SignRequest;
import com.kyriosdata.assinador.domain.SignatureResponse;
import com.kyriosdata.assinador.domain.ValidateRequest;

/**
 * Ponto de entrada do assinador.jar.
 *
 * <p>Uso via linha de comandos:
 * <pre>
 *   java -jar assinador.jar sign    &lt;content&gt; [token]
 *   java -jar assinador.jar validate &lt;content&gt; &lt;signature&gt;
 * </pre>
 *
 * <p>Saída: JSON em stdout com os campos {@code signature}, {@code valid} e {@code message}.
 */
public class Main {

    public static void main(String[] args) {
        if (args.length < 2) {
            printError("Uso: assinador.jar <operação> <content> [token|signature]");
            System.exit(1);
        }

        String operation = args[0];
        String content   = args[1];

        SignatureService service = new FakeSignatureService();

        switch (operation) {
            case "sign" -> {
                String token = args.length > 2 ? args[2] : null;
                SignRequest req = new SignRequest();
                req.setContent(content);
                req.setToken(token);
                SignatureResponse resp = service.sign(req);
                printJson(resp);
                System.exit(resp.isValid() ? 0 : 1);
            }
            case "validate" -> {
                if (args.length < 3) {
                    printError("validate requer: <content> <signature>");
                    System.exit(1);
                }
                String signature = args[2];
                ValidateRequest req = new ValidateRequest();
                req.setContent(content);
                req.setSignature(signature);
                SignatureResponse resp = service.validate(req);
                printJson(resp);
                System.exit(resp.isValid() ? 0 : 1);
            }
            default -> {
                printError("Operação desconhecida: " + operation + ". Use 'sign' ou 'validate'.");
                System.exit(1);
            }
        }
    }

    private static void printJson(SignatureResponse resp) {
        // JSON manual simples, sem dependência externa
        String sig = resp.getSignature() == null ? "null" : "\"" + resp.getSignature() + "\"";
        System.out.printf(
            "{\"signature\":%s,\"valid\":%b,\"message\":\"%s\"}%n",
            sig,
            resp.isValid(),
            resp.getMessage().replace("\"", "\\\"")
        );
    }

    private static void printError(String msg) {
        System.out.printf("{\"signature\":null,\"valid\":false,\"message\":\"%s\"}%n", msg);
    }
}
