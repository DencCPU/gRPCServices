package spothandlers_test

// func TestHandlers_ViewMarket(t *testing.T) {
// 	tDir := t.TempDir()
// 	path := filepath.Join(tDir, "market.txt")
// 	filedate := "Binance, TradingView, Coinbase Exchange, Upbit"
// 	err := os.WriteFile(path, []byte(filedate), 0644)
// 	if err != nil {
// 		t.Fatal("Ошибка записи в файл:", err)
// 	}
// 	logger, _ := logger.NewLogger()
// 	tracerProvider, _ := opentelemetry.NewTracerProviderGrpc(context.Background(), "", "localhost", "4317", 30)
// 	tracer := tracerProvider.Tracer("SpotService")
// 	s, err := memory.NewStorage(logger)
// 	if err != nil {
// 		log.Fatal("ошибка создания хранилища:", err)
// 	}

// 	err = s.AddMarkets(path)
// 	if err != nil {
// 		t.Fatal("ошибка добавдения рынков")
// 	}

// 	re := regexp.MustCompile(`([a-zA-Z0-9]+)+\s?([a-zA-Z0-9]+)?`)
// 	markets := re.FindAllString(filedate, -1)

// 	var checkMap = make(map[string]struct{})
// 	for _, m := range markets {
// 		checkMap[m] = struct{}{}
// 	}

// 	resp := &spot.ViewMarketResp{}
// 	resp.EnableMarkets = make([]*spot.Market, 0, len(checkMap))

// 	for key := range checkMap {
// 		market := spot.Market{MarketName: key}
// 		resp.EnableMarkets = append(resp.EnableMarkets, &market)
// 	}

// 	tests := []struct {
// 		name string // description of this test case
// 		// Named input parameters for receiver constructor.
// 		spotService *usecase.SpotService
// 		// Named input parameters for target function.
// 		req     *spot.ViewMarketReq
// 		want    *spot.ViewMarketResp
// 		wantErr bool
// 	}{
// 		{
// 			name:        "get enable markets",
// 			spotService: usecase.NewSpotInstrument(s, logger, tracer),
// 			req:         &spot.ViewMarketReq{UserRoles: common.UserRole_USER_ROLE_BASIC_USER},
// 			want:        resp,
// 			wantErr:     false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := spothandlers.NewHandlers(tt.spotService)
// 			got, gotErr := h.ViewMarket(context.Background(), tt.req)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("ViewMarket() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("ViewMarket() succeeded unexpectedly")
// 			}
// 			for _, got_el := range got.EnableMarkets {
// 				if _, exists := checkMap[got_el.MarketName]; !exists {
// 					t.Fatal("ошибка. Элемент отсутсвует.")
// 				}
// 			}
// 		})
// 	}
// }
