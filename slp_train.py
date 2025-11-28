import numpy as np
import argparse
import time

# 1. Parse Args
parser = argparse.ArgumentParser()
parser.add_argument('--lr', type=float, default=0.01, help='Learning Rate')
parser.add_argument('--epochs', type=int, default=10, help='Number of epochs')
parser.add_argument('--inputs', type=int, default=3, help='Number of inputs')
args = parser.parse_args()

print(f"Initialize SLP: lr={args.lr}, epochs={args.epochs}")

# 2. Dummy Data (AND/OR gate logic simulation)
# Random input [Samples, Features]
X = np.random.rand(100, args.inputs)
# Random target [Samples, 1]
y = np.random.randint(0, 2, (100, 1))

# 3. Initialize Weights
np.random.seed(42)
weights = np.random.rand(args.inputs, 1)
bias = np.random.rand(1)

# 4. Activation Function (Sigmoid)
def sigmoid(x):
    return 1 / (1 + np.exp(-x))

def sigmoid_derivative(x):
    return x * (1 - x)

# 5. Training Loop
for epoch in range(args.epochs):
    # Forward
    inputs = X
    weighted_sum = np.dot(inputs, weights) + bias
    output = sigmoid(weighted_sum)

    # Error
    error = y - output
    
    # Backprop (Gradient Descent)
    adjustments = error * sigmoid_derivative(output)
    weights += np.dot(inputs.T, adjustments) * args.lr
    bias += np.sum(adjustments) * args.lr
    
    # Log progress every few epochs or last epoch
    if epoch % 2 == 0 or epoch == args.epochs - 1:
        loss = np.mean(np.square(error))
        print(f"Epoch {epoch+1}/{args.epochs} - Loss: {loss:.4f}")
    
    # Simulate some computation time
    time.sleep(0.2)

print("Training Complete.")
