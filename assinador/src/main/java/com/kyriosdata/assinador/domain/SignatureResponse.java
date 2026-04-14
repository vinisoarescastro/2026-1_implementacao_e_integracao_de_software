package com.kyriosdata.assinador.domain;

/**
 * Resposta unificada para operações de assinatura e validação.
 */
public class SignatureResponse {

    private final String signature;
    private final boolean valid;
    private final String message;

    public SignatureResponse(String signature, boolean valid, String message) {
        this.signature = signature;
        this.valid     = valid;
        this.message   = message;
    }

    public String  getSignature() { return signature; }
    public boolean isValid()      { return valid;     }
    public String  getMessage()   { return message;   }
}
