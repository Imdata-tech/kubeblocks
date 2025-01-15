import tushare as ts
import pandas as pd
import os
from datetime import datetime
from ..config import TUSHARE_TOKEN, START_DATE, END_DATE

class DataFetcher:
    def __init__(self):
        ts.set_token(TUSHARE_TOKEN)
        self.pro = ts.pro_api()
        
    def get_hs300_stocks(self):
        """获取沪深300成分股"""
        df = self.pro.index_weight(index_code='000300.SH', 
                                 start_date=START_DATE,
                                 end_date=END_DATE)
        return df['con_code'].unique().tolist()
    
    def get_stock_data(self, ts_code):
        """获取单只股票的历史数据"""
        # Only fetch necessary columns and limit to last 2 years
        fields = ['ts_code', 'trade_date', 'open', 'high', 'low', 'close', 'vol']
        df = ts.pro_bar(ts_code=ts_code,
                       adj='qfq',
                       start_date=(datetime.now().year - 2).strftime('%Y%m%d'),
                       end_date=END_DATE,
                       freq='W',
                       fields=fields)
        return df[fields]  # Ensure only requested columns are returned
    
    def save_all_data(self):
        """获取并保存所有股票数据"""
        stocks = self.get_hs300_stocks()
        os.makedirs('data/historical', exist_ok=True)
        
        success_count = 0
        error_count = 0
        
        # Process in smaller batches with memory cleanup
        batch_size = 20
        for i in range(0, len(stocks), batch_size):
            batch = stocks[i:i + batch_size]
            logger.info(f"Processing batch {i//batch_size + 1}/{(len(stocks)//batch_size)+1}")
            
            for stock in batch:
                try:
                    df = self.get_stock_data(stock)
                    if df is not None and not df.empty:
                        df.to_csv(f'data/historical/{stock}.csv', index=False)
                        success_count += 1
                        logger.info(f"Successfully saved data for {stock}")
                    else:
                        logger.warning(f"No data returned for {stock}")
                except Exception as e:
                    error_count += 1
                    logger.error(f'Error fetching data for {stock}: {str(e)}')
                finally:
                    # Clean up memory
                    del df
                    pd.DataFrame().empty  # Force garbage collection
            
            # Clean up between batches
            pd.DataFrame().empty  # Force garbage collection
        
        logger.info(f"Data fetch completed. Success: {success_count}, Errors: {error_count}")
        return success_count > 0
