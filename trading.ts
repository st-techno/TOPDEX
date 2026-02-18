import { Server } from '@grpc/grpc-js';
import Redis from 'ioredis';

const redis = new Redis();

app.post('/orders', async (req, res) => {
  const { symbol, side, price, size } = req.body;
  // Risk pre-check
  const marginReq = calculateInitialMargin(price, size, leverage);
  const availMargin = await redis.get(`user:margin:${req.user.id}`);
  if (parseFloat(availMargin) < marginReq) throw new Error('Insufficient margin');

  const orderId = await grpcClient.matchingEngine.addOrder({ ...req.body });
  // WS broadcast
  pubsub.publish(symbol, { type: 'order_ack', orderId });
  res.json({ orderId });
});

function calculateFundingRate(markPrice: number, indexPrice: number): number {
  return (markPrice - indexPrice) / indexPrice * 0.01; // Hourly, clamped
}
