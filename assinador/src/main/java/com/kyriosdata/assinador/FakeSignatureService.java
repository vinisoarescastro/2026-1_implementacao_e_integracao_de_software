package com.kyriosdata.assinador;

import com.kyriosdata.assinador.domain.SignRequest;
import com.kyriosdata.assinador.domain.SignatureResponse;
import com.kyriosdata.assinador.domain.ValidateRequest;

/**
 * Implementação simulada de {@link SignatureService}.
 *
 * <p>Retorna uma assinatura pré-construída quando os parâmetros são válidos.
 * O foco está na validação rigorosa dos parâmetros de entrada, conforme US-02.2 e US-02.3.
 */
public class FakeSignatureService implements SignatureService {

    static final String FAKE_SIGNATURE = "MOCKED_SIGNATURE_BASE64_==";

    @Override
    public SignatureResponse sign(SignRequest request) {
        if (request == null) {
            return error("Requisição inválida: objeto nulo");
        }
        if (request.getContent() == null || request.getContent().isBlank()) {
            return error("Parâmetro 'content' inválido ou ausente");
        }
        return new SignatureResponse(FAKE_SIGNATURE, true, "Assinatura criada com sucesso");
    }

    @Override
    public SignatureResponse validate(ValidateRequest request) {
        if (request == null) {
            return error("Requisição inválida: objeto nulo");
        }
        if (request.getContent() == null || request.getContent().isBlank()) {
            return error("Parâmetro 'content' inválido ou ausente");
        }
        if (request.getSignature() == null || request.getSignature().isBlank()) {
            return error("Parâmetro 'signature' inválido ou ausente");
        }

        boolean valid = FAKE_SIGNATURE.equals(request.getSignature());
        String msg = valid ? "Assinatura é válida" : "Assinatura é inválida";
        return new SignatureResponse(request.getSignature(), valid, msg);
    }

    private SignatureResponse error(String message) {
        return new SignatureResponse(null, false, message);
    }
}
