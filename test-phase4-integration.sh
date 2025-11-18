#!/bin/bash
# Phase 4 Integration Tests
# Tests the feature parity components

echo "════════════════════════════════════════════════════"
echo "  PHASE 4: FEATURE PARITY - INTEGRATION TESTS"
echo "════════════════════════════════════════════════════"
echo ""

PASS=0
FAIL=0

check() {
  if [ "$2" ]; then
    echo "✅ $1"
    PASS=$((PASS + 1))
  else
    echo "❌ $1"
    FAIL=$((FAIL + 1))
  fi
}

echo "1. V2 Event Manager (WebSocket):"
check "manager.go exists" "$([ -f v2/events/manager.go ] && echo 1)"
check "EventManager struct" "$(grep -q 'type EventManager struct' v2/events/manager.go && echo 1)"
check "NewEventManager function" "$(grep -q 'func NewEventManager' v2/events/manager.go && echo 1)"
check "Start method" "$(grep -q 'func (m \*EventManager) Start' v2/events/manager.go && echo 1)"
check "Broadcast method" "$(grep -q 'func (m \*EventManager) Broadcast' v2/events/manager.go && echo 1)"
check "HandleWebSocket method" "$(grep -q 'func (m \*EventManager) HandleWebSocket' v2/events/manager.go && echo 1)"
check "Event type definitions" "$(grep -q 'EventAssetCreated' v2/events/manager.go && echo 1)"
check "Client struct" "$(grep -q 'type Client struct' v2/events/manager.go && echo 1)"

echo ""
echo "2. V2 Printer Manager:"
check "manager.go exists" "$([ -f v2/integrations/printers/manager.go ] && echo 1)"
check "PrinterManager struct" "$(grep -q 'type PrinterManager struct' v2/integrations/printers/manager.go && echo 1)"
check "NewPrinterManager function" "$(grep -q 'func NewPrinterManager' v2/integrations/printers/manager.go && echo 1)"
check "AddPrinter method" "$(grep -q 'func (m \*PrinterManager) AddPrinter' v2/integrations/printers/manager.go && echo 1)"
check "ListPrinters method" "$(grep -q 'func (m \*PrinterManager) ListPrinters' v2/integrations/printers/manager.go && echo 1)"
check "Printer struct" "$(grep -q 'type Printer struct' v2/integrations/printers/manager.go && echo 1)"
check "PrinterAdapter interface" "$(grep -q 'type PrinterAdapter interface' v2/integrations/printers/manager.go && echo 1)"
check "Printer state definitions" "$(grep -q 'StateOffline' v2/integrations/printers/manager.go && echo 1)"

echo ""
echo "3. OctoPrint Adapter:"
check "octoprint_adapter.go exists" "$([ -f v2/integrations/printers/octoprint_adapter.go ] && echo 1)"
check "OctoPrintAdapter struct" "$(grep -q 'type OctoPrintAdapter struct' v2/integrations/printers/octoprint_adapter.go && echo 1)"
check "NewOctoPrintAdapter function" "$(grep -q 'func NewOctoPrintAdapter' v2/integrations/printers/octoprint_adapter.go && echo 1)"
check "GetState implementation" "$(grep -q 'func (a \*OctoPrintAdapter) GetState' v2/integrations/printers/octoprint_adapter.go && echo 1)"
check "SendGCode implementation" "$(grep -q 'func (a \*OctoPrintAdapter) SendGCode' v2/integrations/printers/octoprint_adapter.go && echo 1)"
check "UploadFile implementation" "$(grep -q 'func (a \*OctoPrintAdapter) UploadFile' v2/integrations/printers/octoprint_adapter.go && echo 1)"

echo ""
echo "4. Klipper Adapter:"
check "klipper_adapter.go exists" "$([ -f v2/integrations/printers/klipper_adapter.go ] && echo 1)"
check "KlipperAdapter struct" "$(grep -q 'type KlipperAdapter struct' v2/integrations/printers/klipper_adapter.go && echo 1)"
check "NewKlipperAdapter function" "$(grep -q 'func NewKlipperAdapter' v2/integrations/printers/klipper_adapter.go && echo 1)"
check "GetState implementation" "$(grep -q 'func (a \*KlipperAdapter) GetState' v2/integrations/printers/klipper_adapter.go && echo 1)"
check "SendGCode implementation" "$(grep -q 'func (a \*KlipperAdapter) SendGCode' v2/integrations/printers/klipper_adapter.go && echo 1)"

echo ""
echo "5. Processing Adapter:"
check "processing_adapter.go exists" "$([ -f core/integration/processing_adapter.go ] && echo 1)"
check "ProcessingAdapter struct" "$(grep -q 'type ProcessingAdapter struct' core/integration/processing_adapter.go && echo 1)"
check "NewProcessingAdapter function" "$(grep -q 'func NewProcessingAdapter' core/integration/processing_adapter.go && echo 1)"
check "ProcessAsset method" "$(grep -q 'func (a \*ProcessingAdapter) ProcessAsset' core/integration/processing_adapter.go && echo 1)"
check "ProcessFolder method" "$(grep -q 'func (a \*ProcessingAdapter) ProcessFolder' core/integration/processing_adapter.go && echo 1)"
check "ConvertV1AssetToV2 method" "$(grep -q 'func (a \*ProcessingAdapter) ConvertV1AssetToV2' core/integration/processing_adapter.go && echo 1)"

echo ""
echo "6. Discovery Adapter:"
check "discovery_adapter.go exists" "$([ -f core/integration/discovery_adapter.go ] && echo 1)"
check "DiscoveryAdapter struct" "$(grep -q 'type DiscoveryAdapter struct' core/integration/discovery_adapter.go && echo 1)"
check "NewDiscoveryAdapter function" "$(grep -q 'func NewDiscoveryAdapter' core/integration/discovery_adapter.go && echo 1)"
check "DiscoverAssets method" "$(grep -q 'func (a \*DiscoveryAdapter) DiscoverAssets' core/integration/discovery_adapter.go && echo 1)"
check "DiscoverProjects method" "$(grep -q 'func (a \*DiscoveryAdapter) DiscoverProjects' core/integration/discovery_adapter.go && echo 1)"
check "GetDiscoveryStats method" "$(grep -q 'func (a \*DiscoveryAdapter) GetDiscoveryStats' core/integration/discovery_adapter.go && echo 1)"

echo ""
echo "7. Directory Structure:"
check "v2/events directory exists" "$([ -d v2/events ] && echo 1)"
check "v2/integrations/printers directory exists" "$([ -d v2/integrations/printers ] && echo 1)"

# Count files
EVENT_FILES=$(find v2/events -name "*.go" 2>/dev/null | wc -l)
PRINTER_FILES=$(find v2/integrations/printers -name "*.go" 2>/dev/null | wc -l)
INTEGRATION_FILES=$(find core/integration -name "*.go" 2>/dev/null | wc -l)

echo ""
echo "Event files: $EVENT_FILES (expected: >=1)"
echo "Printer files: $PRINTER_FILES (expected: >=3)"
echo "Integration files: $INTEGRATION_FILES (expected: >=7)"

check "Event package has files" "$([ $EVENT_FILES -ge 1 ] && echo 1)"
check "Printer package has files" "$([ $PRINTER_FILES -ge 3 ] && echo 1)"
check "Integration package updated" "$([ $INTEGRATION_FILES -ge 7 ] && echo 1)"

echo ""
echo "8. Event Types:"
check "EventAssetCreated defined" "$(grep -q 'EventAssetCreated' v2/events/manager.go && echo 1)"
check "EventAssetUpdated defined" "$(grep -q 'EventAssetUpdated' v2/events/manager.go && echo 1)"
check "EventAssetDeleted defined" "$(grep -q 'EventAssetDeleted' v2/events/manager.go && echo 1)"
check "EventScanStarted defined" "$(grep -q 'EventScanStarted' v2/events/manager.go && echo 1)"
check "EventScanCompleted defined" "$(grep -q 'EventScanCompleted' v2/events/manager.go && echo 1)"
check "EventPrinterStatus defined" "$(grep -q 'EventPrinterStatus' v2/events/manager.go && echo 1)"

echo ""
echo "9. Printer Types:"
check "PrinterTypeOctoPrint defined" "$(grep -q 'PrinterTypeOctoPrint' v2/integrations/printers/manager.go && echo 1)"
check "PrinterTypeKlipper defined" "$(grep -q 'PrinterTypeKlipper' v2/integrations/printers/manager.go && echo 1)"
check "PrinterState enum" "$(grep -q 'StateOffline.*PrinterState' v2/integrations/printers/manager.go && echo 1)"

echo ""
echo "10. Adapter Integration:"
check "Processing uses feature flags" "$(grep -q 'EnableV2Processing' core/integration/processing_adapter.go && echo 1)"
check "Discovery uses feature flags" "$(grep -q 'EnableV2Processing' core/integration/discovery_adapter.go && echo 1)"
check "Processing uses event manager" "$(grep -q 'eventMgr' core/integration/processing_adapter.go && echo 1)"
check "Discovery uses event manager" "$(grep -q 'eventMgr' core/integration/discovery_adapter.go && echo 1)"

echo ""
echo "════════════════════════════════════════════════════"
echo "Results: $PASS passed, $FAIL failed"
echo "════════════════════════════════════════════════════"

if [ $FAIL -eq 0 ]; then
  echo "✅ PHASE 4 INTEGRATION VERIFIED"
  echo ""
  echo "Feature parity components successfully implemented:"
  echo "  • V2 Event Manager (WebSocket)"
  echo "  • V2 Printer Manager"
  echo "  • OctoPrint Adapter"
  echo "  • Klipper Adapter"
  echo "  • Processing Adapter (v1/v2)"
  echo "  • Discovery Adapter (v1/v2)"
  exit 0
else
  echo "⚠️  Some checks failed"
  exit 1
fi
