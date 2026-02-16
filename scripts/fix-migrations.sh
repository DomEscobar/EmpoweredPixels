#!/bin/bash
DIR="/root/EmpoweredPixels/backend/internal/infra/db/migrations"
for f in "$DIR"/*.sql; do
    # 1. Fix double "IF NOT EXISTS"
    sed -i 's/IF NOT EXISTS IF NOT EXISTS/IF NOT EXISTS/g' "$f"
    
    # 2. Make ALTER TABLE ... ADD COLUMN idempotent
    # Convert 'ALTER TABLE table_name ADD COLUMN column_name ...;' 
    # To 'DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='\''table_name'\'' AND column_name='\''column_name'\'') THEN ALTER TABLE table_name ADD COLUMN column_name ...; END IF; END $$;'
    
    # Simple regex approach for common patterns in this project
    perl -i -pe "s/ALTER TABLE (\w+) ADD COLUMN (\w+) ([^;]+);/DO \\\$\\\$ BEGIN IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='\1' AND column_name='\2') THEN ALTER TABLE \1 ADD COLUMN \2 \3; END IF; END \\\$\\\$ /g" "$f"
done
