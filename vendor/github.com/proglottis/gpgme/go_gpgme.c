#include "go_gpgme.h"

unsigned int key_revoked(gpgme_key_t k) {
	return k->revoked;
}

unsigned int key_expired(gpgme_key_t k) {
	return k->expired;
}

unsigned int key_disabled(gpgme_key_t k) {
	return k->disabled;
}

unsigned int key_invalid(gpgme_key_t k) {
	return k->invalid;
}

unsigned int key_can_encrypt(gpgme_key_t k) {
	return k->can_encrypt;
}

unsigned int key_can_sign(gpgme_key_t k) {
	return k->can_sign;
}

unsigned int key_can_certify(gpgme_key_t k) {
	return k->can_certify;
}

unsigned int key_secret(gpgme_key_t k) {
	return k->secret;
}

unsigned int key_can_authenticate(gpgme_key_t k) {
	return k->can_authenticate;
}

unsigned int key_is_qualified(gpgme_key_t k) {
	return k->is_qualified;
}

unsigned int signature_wrong_key_usage(gpgme_signature_t s) {
    return s->wrong_key_usage;
}

unsigned int signature_pka_trust(gpgme_signature_t s) {
    return s->pka_trust;
}

unsigned int signature_chain_model(gpgme_signature_t s) {
    return s->chain_model;
}

extern unsigned int subkey_revoked(gpgme_subkey_t k) {
	return k->revoked;
}

extern unsigned int subkey_expired(gpgme_subkey_t k) {
	return k->expired;
}

extern unsigned int subkey_disabled(gpgme_subkey_t k) {
	return k->disabled;
}

extern unsigned int subkey_invalid(gpgme_subkey_t k) {
	return k->invalid;
}

extern unsigned int subkey_secret(gpgme_subkey_t k) {
	return k->secret;
}

unsigned int uid_revoked(gpgme_user_id_t u) {
	return u->revoked;
}

unsigned int uid_invalid(gpgme_user_id_t u) {
	return u->invalid;
}
