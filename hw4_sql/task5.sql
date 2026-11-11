SELECT item, price
FROM orders
WHERE price = (
    SELECT MIN(price)
    FROM orders
    WHERE price > 10000
);