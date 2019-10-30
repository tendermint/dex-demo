import {
  ChartingLibraryWidgetOptions,
  IBasicDataFeed,
  LibrarySymbolInfo
} from "../../public/charting_library/charting_library.min";
import store from "../store";
import chartDataProvider, {formatTVCandles, serializeUID} from "./ChartDataProvider";

export const TVDatafeed: IBasicDataFeed = {

  onReady: cb => {
    // TradingView expects asynchronous execution of this callback
    setTimeout(() => cb({
      supported_resolutions: ["1", "5", "15", "60"],
    }));
  },

  // only need searchSymbols when search is enabled
  searchSymbols() {},

  resolveSymbol: (symbolName, onSymbolResolvedCallback, onResolveErrorCallback) => {
    const { assets: { assets, symbolToAssetId } } = store.getState();
    const [ _, quoteSymbol ] = symbolName.split('/');
    const quoteAssetId = symbolToAssetId[quoteSymbol];
    const quoteAsset = assets[quoteAssetId];

    const symbol_stub: LibrarySymbolInfo = {
      name: symbolName,
      full_name: '',
      exchange: '',
      listed_exchange: '',
      description: '',
      type: 'crypto',
      session: '24x7',
      // timezone doesn't really matter here since the market is 24X7
      timezone: 'America/Los_Angeles',
      ticker: symbolName,
      minmov: 1,
      pricescale: 10 ** Math.min(4, quoteAsset.nativeDecimals),
      has_intraday: true,
      intraday_multipliers: ['1'],
      supported_resolutions:  ["1", "5", "15", "60"],
      volume_precision: Math.min(4, quoteAsset.nativeDecimals),
      data_status: 'streaming',
    };

    // TradingView expects asynchronous execution of this callback
    setTimeout(() => onSymbolResolvedCallback(symbol_stub), 0);
  },

  getBars: async (symbolInfo, resolution, from, to, onHistoryCallback, onErrorCallback, firstDataRequest) => {
    const {
      assets: { assets, symbolToAssetId },
      exchange: { pairToMarketId },
    } = store.getState();

    const [ baseSymbol, quoteSymbol ] = symbolInfo.name.split('/');
    const baseAssetId = symbolToAssetId[baseSymbol];
    const baseAsset = assets[baseAssetId];
    const quoteAssetId = symbolToAssetId[quoteSymbol];
    const quoteAsset = assets[quoteAssetId];

    const marketId = pairToMarketId[`${baseAsset.symbol}/${quoteAsset.symbol}`]

    if (!baseAsset) return onErrorCallback(`Cannot get bars for ${baseSymbol}`);
    if (!quoteAsset) return onErrorCallback(`Cannot get bars for ${quoteSymbol}`);

    try {
      const rawCandles = await chartDataProvider.fetchCandles(
        marketId,
        baseSymbol,
        quoteSymbol,
        resolution,
        2000,
        to * 1000,
        from * 1000,
      );
      onHistoryCallback(formatTVCandles(rawCandles, 1 / (10 ** quoteAsset.decimals)), { noData: !rawCandles.length });
    } catch (e) {
      onErrorCallback(e.message);
    }

  },

  subscribeBars: (symbolInfo, resolution, onRealtimeCallback, subscriberUID, onResetCacheNeededCallback) => {
    const uid = serializeUID(symbolInfo, resolution);
    chartDataProvider.subscribe(uid, onRealtimeCallback);
  },

  unsubscribeBars: (subscriberUID) => {
    chartDataProvider.unsubscribe(subscriberUID);
  },

};

/* optional methods */
// getServerTime: cb => {},
// calculateHistoryDepth: (resolution, resolutionBack, intervalBack) => {},
// getMarks: (symbolInfo, startDate, endDate, onDataCallback, resolution) => {},
// getTimeScaleMarks: (symbolInfo, startDate, endDate, onDataCallback, resolution) => {},

export const baseWidgetOption: ChartingLibraryWidgetOptions = {
  symbol: '',
  interval: '',
  debug: false,
  datafeed: TVDatafeed,
  container_id: 'tv-container',
  library_path: '/charting_library/',
  locale: 'en',
  disabled_features: ['use_localstorage_for_settings', 'header_widget', 'compare_symbol', 'symbol_search_hot_key'],
  enabled_features: ['side_toolbar_in_fullscreen_mode'],
  fullscreen: false,
  autosize: true,
  theme: 'Dark',
  custom_css_url: './custom.css',
  overrides: {
    'mainSeriesProperties.style': 1,
    'paneProperties.topMargin': 10,
    'paneProperties.bottomMargin': 8,
    'paneProperties.background': '#101010',
    'paneProperties.gridProperties.color': "#2c2c2c",
    'paneProperties.vertGridProperties.color': '#2c2c2c',
    'paneProperties.horzGridProperties.color': '#2c2c2c',
    'scalesProperties.lineColor': '#3c3c3c',
    'scalesProperties.textColor': '#8a8a8a',
  },
}
