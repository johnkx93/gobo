Common built-in validation tags (useful quick reference)
required — must be non-zero.
required_with, required_without, required_if — conditional required.
email — email format (simple/common check).
url — valid URL.
min, max — numeric or string length (e.g., min=3, max=50).
len — exact length (strings/arrays).
gt, gte, lt, lte — greater/less comparisons.
oneof — value must be one of listed (e.g., oneof=admin user guest).
eqfield — equal to another field (e.g., password confirmation).
alpha, alphanum, alphanumunicode — char class checks.
numeric, number — numeric checks.
uuid, uuid4 — UUID format checks.
isbn, isbn10, isbn13, creditcard, base64 — other specialized checks.
(There's a long list; see the validator docs for completeness.)