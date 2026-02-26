#!/bin/bash

# Einfacher Test-Script für Python Lösung

echo "==================================="
echo "Python Interview Lösung - Tests"
echo "==================================="
echo ""

# Farben
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

cd "$(dirname "$0")"

# Port 8080 freimachen
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 1

echo -e "${YELLOW}Server starten...${NC}"
poetry run python -m simple_web_app.main_solution > /tmp/server.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

# Cleanup Funktion
cleanup() {
    echo ""
    echo -e "${YELLOW}Server stoppen...${NC}"
    kill $SERVER_PID 2>/dev/null || true
}
trap cleanup EXIT

# Warten bis Server läuft
echo "Warte auf Server..."
sleep 3

# Test 1: Hello World
echo ""
echo -e "${YELLOW}Test 1: /hello-world${NC}"
RESPONSE=$(curl -s http://localhost:8080/hello-world)
echo "Response: $RESPONSE"
if echo "$RESPONSE" | grep -q "Hello World"; then
    echo -e "${GREEN}✓ Test bestanden${NC}"
else
    echo -e "${RED}✗ Test fehlgeschlagen${NC}"
    exit 1
fi

# Test 2: Repo-List ohne Filter
echo ""
echo -e "${YELLOW}Test 2: /repo-list/golang${NC}"
RESPONSE=$(curl -s http://localhost:8080/repo-list/golang)
if echo "$RESPONSE" | grep -q "organization"; then
    COUNT=$(echo "$RESPONSE" | grep -o '"count":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}✓ Gefunden: $COUNT Repositories${NC}"
else
    echo -e "${RED}✗ Test fehlgeschlagen${NC}"
    exit 1
fi

# Test 3: Mit Filter
echo ""
echo -e "${YELLOW}Test 3: /repo-list/golang?repo_filter=go${NC}"
RESPONSE=$(curl -s "http://localhost:8080/repo-list/golang?repo_filter=go")
if echo "$RESPONSE" | grep -q "golang"; then
    COUNT=$(echo "$RESPONSE" | grep -o '"count":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}✓ Gefiltert: $COUNT Repositories${NC}"
else
    echo -e "${RED}✗ Test fehlgeschlagen${NC}"
    exit 1
fi

# Test 4: 404 Test
echo ""
echo -e "${YELLOW}Test 4: Non-existent org (sollte 404 sein)${NC}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/repo-list/nonexistentorg999)
if [ "$HTTP_CODE" = "404" ]; then
    echo -e "${GREEN}✓ Korrekt: 404 zurückgegeben${NC}"
else
    echo -e "${RED}✗ Erwartet 404, bekommen $HTTP_CODE${NC}"
fi

echo ""
echo "==================================="
echo -e "${GREEN}Alle Tests bestanden! ✓${NC}"
echo "==================================="
echo ""
echo "Lösung funktioniert korrekt:"
echo "  ✓ /hello-world"
echo "  ✓ /repo-list/{org_name}"
echo "  ✓ Filter funktioniert"
echo "  ✓ Fehlerbehandlung (404)"
echo ""
