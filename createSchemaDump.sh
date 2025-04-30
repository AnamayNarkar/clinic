#!/bin/bash
pg_dump -U postgres -h localhost -p 5432 -s -F p -v -f sql/schema/clinic_dump.sql clinic
