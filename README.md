# Payment API

The payment API receives Stripe webhooks and publishes payment events to RabbitMQ. Stripe webhook signatures are verified with `STRIPE_SECRET_ENDPOINT`; events with an invalid signature are rejected.

## Barber appointment events

Barber events are published to the `barber-api` RabbitMQ topic only when the Stripe PaymentIntent metadata identifies a Barber payment and contains a valid appointment UUID:

| Metadata key | Required value |
| --- | --- |
| `payment_context` | `barber` |
| `appointment_id` | Appointment UUID used by the Barber API |

Metadata without the explicit Barber context, without an appointment ID, or with a malformed UUID does not produce a Barber event. Keep these values on the PaymentIntent. For refunds, the webhook uses an expanded PaymentIntent from the charge when available, otherwise retrieves the PaymentIntent by ID.

The published message uses the standard message-worker envelope (`event` and `data`). Its `data` object is normalized as follows:

| Stripe event | RabbitMQ event | `data` fields | Appointment status |
| --- | --- | --- | --- |
| `payment_intent.succeeded` | `appointment_status_update` | `appointment_id`, `status`, `payment_event`, `payment_intent_id`, `charge_id` | `scheduled` |
| `charge.refunded` | `appointment_status_update` | `appointment_id`, `status`, `payment_event`, `payment_intent_id`, `charge_id` | `canceled` |

`appointment_status_update` is the event name already consumed by the Barber API. `payment_event` preserves which Stripe event caused the update. `status` is the resulting Barber appointment status. A refund event marks the appointment `canceled`; this contract applies to any Stripe refund amount reported by `charge.refunded`. Stripe events that do not contain valid Barber metadata do not publish an appointment status update.

## Existing payment events

The successful-payment transaction workflow remains in place, including the existing `akademia-api` / `user.new_course` publication based on course metadata. Subscription and customer webhook handling is unchanged. Unknown Stripe event types continue to be logged and ignored.
