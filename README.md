# Payment API

The Payment API creates Stripe products and prices from the Barber service catalog, starts Stripe Checkout for appointments, records signed Stripe webhook results, and reports revenue. All money values in API payloads are integer minor units (for example, cents for USD).

## Barber service catalog contract

Barber App publishes `service.created` to the `finance-api` RabbitMQ topic after saving a service. The event data contains:

| Field | Contract |
| --- | --- |
| `service_id` | Barber service UUID |
| `name` | Service display name |
| `description` | Service description |
| `price` | Positive integer in the currency's minor unit |
| `currency` | ISO currency code, such as `usd` |

Payment API persists a trusted catalog record keyed by `barber_service_id`, creates its Stripe product and price, then publishes `service.payment_details_created` to the `barber-app` topic with `service_id`, `external_id` (Stripe product ID), and `external_price_id` (Stripe Price ID). Repeated events update the local trusted record and reuse its Stripe product. Checkout uses this server-side record; browser-provided price IDs and totals are not accepted.

## Appointment Checkout

`POST /checkout/appointments` starts a hosted Stripe Checkout Session. It is a server-to-server endpoint and requires `Authorization: Bearer <BARBER_PAYMENT_API_TOKEN>`; configure the same secret in the Barber service and Payment API. The request accepts identifiers only:

```json
{
  "appointment_id": "<appointment-uuid>",
  "service_ids": ["<service-uuid-1>", "<service-uuid-2>"]
}
```

The Payment API validates UUIDs, rejects repeated or unavailable services and mixed currencies, looks up each service's current price and Stripe Price ID in its synchronized catalog, and calculates the total from those records. Browser-supplied totals, durations, prices, Stripe Price IDs, and roles are ignored. Barber API must validate the appointment, its selected service IDs, service availability, offered services, and scheduling constraints before calling this internal endpoint.

Successful response (`201 Created`):

```json
{
  "checkout_url": "https://checkout.stripe.com/...",
  "session_id": "cs_...",
  "amount_total": 6500,
  "currency": "usd",
  "payment_status": "pending"
}
```

The API stores a pending payment record before returning the URL. Repeating the same request while that Checkout Session remains pending returns its existing URL. A different service selection while a payment is pending, or a new checkout after the appointment is paid/refunded, returns `409 Conflict`. Configure `STRIPE_SUCCESS_URL` and `STRIPE_CANCEL_URL` for the hosted Checkout redirects.

## Separate appointment and payment statuses

Appointment status belongs to Barber API and is never changed by a Payment API event. Appointment statuses are `pending`, `accepted`, `refused`, and `canceled`: Barber creates an appointment as `pending`; a barber's accept/refuse action moves it to `accepted` or `refused`; a customer or barber cancellation action moves it to `canceled`. Payment API does not perform those actions.

Payment status is separate: `unpaid`, `pending`, `paid`, or `refunded`.

| Action | Payment status change | Appointment status change |
| --- | --- | --- |
| Appointment has no checkout, or a signed failed/expired Checkout event arrives | `unpaid` | None |
| Hosted Checkout Session is created | `pending` | None |
| Verified `payment_intent.succeeded` webhook | `paid` | None |
| Verified `charge.refunded` webhook | `refunded` | None |
| Barber accepts, refuses, or cancels an appointment | None | Barber API changes the appointment status |

Refunds only update Payment API records and the payment status. They do not cancel appointments or trigger email/SMS notifications. Refund workflows and customer notifications remain future CRM work.

## Barber payment-status event contract

This updates the previous `appointment_status_update` contract: Payment API no longer sends appointment `status` values such as `scheduled` or `canceled`. It publishes `appointment_payment_status_update` to the `barber-app` topic with an `appointment_id` and `payment_status`; it may also include `session_id`, `payment_intent_id`, `charge_id`, `amount_total`, and `currency`. Barber consumers must update only the payment-status field and must not map this event to appointment status.

Events are sent when checkout becomes pending, when a signed payment-failure/expiry event makes it unpaid, after successful payment, and after a refund. The payment callback includes `payment_status` rather than `status` so appointment lifecycle state remains independent.

## Stripe webhooks and payment records

`POST /webhook/stripe` continues to verify `Stripe-Signature` with `STRIPE_SECRET_ENDPOINT`. Do not disable signature or timestamp validation. `STRIPE_IGNORE_API_VERSION_MISMATCH=true` only permits local Stripe CLI API-version mismatch handling.

Appointment transactions are keyed idempotently by Stripe Checkout Session, checkout attempt, and PaymentIntent IDs. Repeated success events update the same record. `charge.refunded` updates the charge's cumulative refunded amount and payment status rather than creating another revenue entry. Payment records store the appointment ID, gross amount, currency, Stripe receipt URL, actual Stripe fee when Stripe's balance transaction is available, refund amount, and net amount. `charge.updated` refreshes fee data when it becomes available.

## Authorized receipts and revenue

Configure `PAYMENT_API_ADMIN_TOKEN`. Transaction routes (including receipt URLs) and admin reports require `Authorization: Bearer <PAYMENT_API_ADMIN_TOKEN>`; requests fail closed if the token is not configured. This also protects the existing admin dashboard endpoint.

`GET /admin/revenue` returns one summary per currency so unlike currencies are never added together:

```json
{
  "revenue": [
    {
      "currency": "usd",
      "gross": 10000,
      "stripe_fees": 320,
      "refunds": 2000,
      "net_company_revenue": 7680,
      "payment_count": 2,
      "stripe_fee_pending_payments": 1
    }
  ]
}
```

`gross` is the sum of captured amounts for paid/refunded transactions. `refunds` is Stripe's cumulative refunded amount. `stripe_fees` sums actual fee amounts recorded from Stripe balance transactions. `net_company_revenue` is `gross - refunds - stripe_fees`. `stripe_fee_pending_payments` counts captured transactions whose actual fee is not available yet; net revenue is provisional for those payments until Stripe supplies the fee.
