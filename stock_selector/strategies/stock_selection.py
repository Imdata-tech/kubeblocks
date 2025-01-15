class StockSelection:
    def __init__(self, df):
        self.df = df
        
    def ma_crossover(self):
        """判断5周线上穿10周线"""
        self.df['ma_cross'] = (self.df['ma5'].shift(1) < self.df['ma10'].shift(1)) & \
                            (self.df['ma5'] > self.df['ma10'])
        return self.df
        
    def volume_crossover(self):
        """判断成交量指标金叉"""
        self.df['volume_cross'] = (self.df['volume_ma5'].shift(1) < self.df['volume_ma10'].shift(1)) & \
                                (self.df['volume_ma5'] > self.df['volume_ma10'])
        return self.df
        
    def macd_crossover(self):
        """判断MACD指标金叉"""
        self.df['macd_cross'] = (self.df['macd'].shift(1) < self.df['macdsignal'].shift(1)) & \
                              (self.df['macd'] > self.df['macdsignal'])
        return self.df
        
    def kdj_crossover(self):
        """判断KDJ指标金叉"""
        self.df['kdj_cross'] = (self.df['K'].shift(1) < self.df['D'].shift(1)) & \
                             (self.df['K'] > self.df['D'])
        return self.df
        
    def select_stocks(self):
        """综合判断所有条件"""
        self.ma_crossover()
        self.volume_crossover()
        self.macd_crossover()
        self.kdj_crossover()
        
        # 筛选满足所有条件的股票
        selected = self.df[
            (self.df['ma_cross'] == True) &
            (self.df['volume_cross'] == True) &
            (self.df['macd_cross'] == True) &
            (self.df['kdj_cross'] == True)
        ]
        return selected
