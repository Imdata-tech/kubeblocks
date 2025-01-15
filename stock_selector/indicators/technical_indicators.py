import pandas as pd
import talib

class TechnicalIndicators:
    def __init__(self, df):
        self.df = df
        
    def calculate_ma(self):
        """计算5周线和10周线"""
        self.df['ma5'] = self.df['close'].rolling(window=5).mean()
        self.df['ma10'] = self.df['close'].rolling(window=10).mean()
        return self.df
        
    def calculate_volume(self):
        """计算成交量指标"""
        self.df['volume_ma5'] = self.df['vol'].rolling(window=5).mean()
        self.df['volume_ma10'] = self.df['vol'].rolling(window=10).mean()
        return self.df
        
    def calculate_macd(self):
        """计算MACD指标"""
        self.df['macd'], self.df['macdsignal'], self.df['macdhist'] = talib.MACD(
            self.df['close'], fastperiod=12, slowperiod=26, signalperiod=9)
        return self.df
        
    def calculate_kdj(self):
        """计算KDJ指标"""
        low_list = self.df['low'].rolling(window=9).min()
        high_list = self.df['high'].rolling(window=9).max()
        rsv = (self.df['close'] - low_list) / (high_list - low_list) * 100
        
        self.df['K'] = rsv.ewm(com=2).mean()
        self.df['D'] = self.df['K'].ewm(com=2).mean()
        self.df['J'] = 3 * self.df['K'] - 2 * self.df['D']
        return self.df
        
    def calculate_all(self, indicators=None):
        """
        计算指定的技术指标
        indicators: 要计算的指标列表，可选值：'ma', 'volume', 'macd', 'kdj'
                    如果为None，则计算所有指标
        """
        if indicators is None:
            indicators = ['ma', 'volume', 'macd', 'kdj']
            
        if 'ma' in indicators:
            self.calculate_ma()
        if 'volume' in indicators:
            self.calculate_volume()
        if 'macd' in indicators:
            self.calculate_macd()
        if 'kdj' in indicators:
            self.calculate_kdj()
            
        # 清理中间计算列
        self.cleanup_intermediate_columns()
        return self.df
        
    def cleanup_intermediate_columns(self):
        """清理中间计算列"""
        # 清理KDJ计算中间列
        if 'K' in self.df.columns and 'D' in self.df.columns:
            self.df.drop(columns=['K', 'D'], inplace=True)
