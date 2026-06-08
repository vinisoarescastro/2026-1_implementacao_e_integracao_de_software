package com.kyriosdata.assinador.service;

import com.kyriosdata.assinador.domain.SignRequest;
import com.kyriosdata.assinador.domain.SignatureResponse;
import com.kyriosdata.assinador.domain.ValidateRequest;

/**
 * Interface principal do serviço de assinatura digital.
 * A implementação real usaria PKCS#11 (Sprint 3).
 * Por ora, {@link FakeSignatureService} provê a simulação.
 */
public interface SignatureService {
    SignatureResponse sign(SignRequest request);
    SignatureResponse validate(ValidateRequest request);
}
