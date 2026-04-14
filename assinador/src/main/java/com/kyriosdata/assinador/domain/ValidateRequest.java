package com.kyriosdata.assinador.domain;

/** DTO para requisição de validação de assinatura. */
public class ValidateRequest {
    private String content;
    private String signature;

    public ValidateRequest() {}

    public String getContent() { return content; }
    public void setContent(String content) { this.content = content; }

    public String getSignature() { return signature; }
    public void setSignature(String signature) { this.signature = signature; }
}
