package com.kyriosdata.assinador;

import com.kyriosdata.assinador.domain.SignRequest;
import com.kyriosdata.assinador.domain.SignatureResponse;
import com.kyriosdata.assinador.domain.ValidateRequest;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FakeSignatureServiceTest {

    private final SignatureService service = new FakeSignatureService();

    @Test
    void sign_retornaAssinaturaSimuladaParaConteudoValido() {
        SignRequest req = new SignRequest();
        req.setContent("documento.pdf");

        SignatureResponse resp = service.sign(req);

        assertTrue(resp.isValid());
        assertEquals(FakeSignatureService.FAKE_SIGNATURE, resp.getSignature());
    }

    @Test
    void sign_retornaErroParaConteudoNulo() {
        SignRequest req = new SignRequest();
        req.setContent(null);

        SignatureResponse resp = service.sign(req);

        assertFalse(resp.isValid());
        assertNull(resp.getSignature());
        assertTrue(resp.getMessage().contains("content"));
    }

    @Test
    void sign_retornaErroParaConteudoVazio() {
        SignRequest req = new SignRequest();
        req.setContent("   ");

        SignatureResponse resp = service.sign(req);

        assertFalse(resp.isValid());
    }

    @Test
    void validate_retornaValidoParaAssinaturaCorreta() {
        ValidateRequest req = new ValidateRequest();
        req.setContent("documento.pdf");
        req.setSignature(FakeSignatureService.FAKE_SIGNATURE);

        SignatureResponse resp = service.validate(req);

        assertTrue(resp.isValid());
        assertEquals("Assinatura é válida", resp.getMessage());
    }

    @Test
    void validate_retornaInvalidoParaAssinaturaErrada() {
        ValidateRequest req = new ValidateRequest();
        req.setContent("documento.pdf");
        req.setSignature("assinatura-errada");

        SignatureResponse resp = service.validate(req);

        assertFalse(resp.isValid());
        assertEquals("Assinatura é inválida", resp.getMessage());
    }

    @Test
    void validate_retornaErroParaSignatureNula() {
        ValidateRequest req = new ValidateRequest();
        req.setContent("documento.pdf");
        req.setSignature(null);

        SignatureResponse resp = service.validate(req);

        assertFalse(resp.isValid());
        assertTrue(resp.getMessage().contains("signature"));
    }
}
