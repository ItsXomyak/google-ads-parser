#!/bin/bash
curl -X POST http://localhost:8080/domains \
  -H "Content-Type: application/json" \
  -d '[
    {
      "domain": "youtube.com",
      "legal_name": "Google LLC",
      "country": "США"
    }
  ]'

curl http://localhost:8080/domains
