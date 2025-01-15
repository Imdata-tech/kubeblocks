import pandas as pd
import logging
import os
from datetime import datetime
from data.data_fetcher import DataFetcher
from indicators.technical_indicators import TechnicalIndicators
from strategies.stock_selection import StockSelection

# Configure logging
os.makedirs('logs', exist_ok=True)
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    filename='logs/stock_selector.log'
)
logger = logging.getLogger(__name__)

def main():
    # Initialize data fetcher
    fetcher = DataFetcher()
    
    # Fetch and save all stock data
    logger.info("Fetching stock data...")
    os.makedirs('data/historical', exist_ok=True)
    try:
        fetcher.save_all_data()
        logger.info("Successfully fetched and saved stock data")
    except Exception as e:
        logger.error(f'Failed to fetch stock data: {str(e)}')
        raise
    
    # Get HS300 stocks
    stocks = fetcher.get_hs300_stocks()
    
    # Create results directory and file
    os.makedirs('results', exist_ok=True)
    results_file = f"results/results_{datetime.now().strftime('%Y%m%d_%H%M%S')}.csv"
    with open(results_file, 'w') as f:
        f.write("code,date,close\n")
    
    # Process in smaller batches
    batch_size = 20
    for i in range(0, len(stocks), batch_size):
        batch = stocks[i:i + batch_size]
        logger.info(f"Processing batch {i//batch_size + 1}/{(len(stocks)//batch_size)+1}")
        
        for stock in batch:
            try:
                # Read single stock data
                df = pd.read_csv(f'data/historical/{stock}.csv')
                
                # Calculate only necessary technical indicators
                indicator = TechnicalIndicators(df)
                df = indicator.calculate_all(indicators=['ma', 'macd'])  # 只计算MA和MACD
                
                # Execute stock selection strategy
                selector = StockSelection(df)
                selected = selector.select_stocks()
                
                # If stock meets criteria, save to file
                if not selected.empty:
                    with open(results_file, 'a') as f:
                        f.write(f"{stock},{selected.iloc[-1]['trade_date']},{selected.iloc[-1]['close']}\n")
                    
            except Exception as e:
                logger.error(f'Error processing stock {stock}: {str(e)}')
            finally:
                # Clear memory aggressively
                del df
                del indicator
                del selector
                del selected
                pd.DataFrame().empty  # Force garbage collection
        
        # Clear memory between batches
        pd.DataFrame().empty  # Force garbage collection
    
    logger.info(f"Processing complete. Results saved to {results_file}")

if __name__ == '__main__':
    main()
