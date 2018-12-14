// The package network provides data holders and utilities to describe Artificial Neural Network
package network

import (
	"math"
	"fmt"
	"errors"
)

var (
	// The error to be raised when maximal number of network activation attempts exceeded
	NetErrExceededMaxActivationAttempts = errors.New("maximal network activation attempts exceeded.")
	// The error to be raised when unsupported sensors data array size provided
	NetErrUnsupportedSensorsArraySize = errors.New("the sensors array size is unsupported by network solver")
)

// Defines network solver interface which describes neural network structures with methods to run activation waves through
// them.
type NetworkSolver interface {
	// Propagates activation wave through all network nodes provided number of steps in forward direction.
	// Returns true if activation wave passed from all inputs to outputs.
	ForwardSteps(steps int) (bool, error)

	// Propagates activation wave through all network nodes provided number of steps by recursion from output nodes
	// Returns true if activation wave passed from all inputs to outputs.
	RecursiveSteps() (bool, error)

	// Attempts to relax network given amount of steps until giving up. The network considered relaxed when absolute
	// value of the change at any given point is less than maxAllowedSignalDelta during activation waves propagation.
	// If maxAllowedSignalDelta value is less than or equal to 0, the method will return true without checking for relaxation.
	Relax(maxSteps int, maxAllowedSignalDelta float64) (bool, error)

	// Flushes network state by removing all current activations. Returns true if network flushed successfully or
	// false in case of error.
	Flush() (bool, error)

	// Set sensors values to the input nodes of the network
	LoadSensors(inputs []float64) error
	// Read output values from the output nodes of the network
	ReadOutputs() []float64

	// Returns the total number of neural units in the network
	NodeCount() int
	// Returns the total number of links between nodes in the network
	LinkCount() int
}

// NNodeType defines the type of NNode to create
type NodeType byte

// Predefined NNode types
const (
	// The neuron type
	NeuronNode NodeType = iota
	// The sensor type
	SensorNode
)

// Returns human readable NNode type name for given constant value
func NodeTypeName(ntype NodeType) string {
	switch ntype {
	case NeuronNode:
		return "NEURON"
	case SensorNode:
		return "SENSOR"
	default:
		return "!!! UNKNOWN NODE TYPE !!!"
	}
}

// NeuronType defines the type of neuron to create
type NodeNeuronType byte

// These are NNode layer type
const (
	// The node is in hidden layer
	HiddenNeuron NodeNeuronType = iota
	// The node is in input layer
	InputNeuron
	// The node is in output layer
	OutputNeuron
	// The node is bias
	BiasNeuron
)

// Returns human readable neuron type name for given constant
func NeuronTypeName(nlayer NodeNeuronType) string {
	switch nlayer {
	case HiddenNeuron:
		return "HIDN"
	case InputNeuron:
		return "INPT"
	case OutputNeuron:
		return "OUTP"
	case BiasNeuron:
		return "BIAS"
	default:
		return "!!! UNKNOWN NEURON TYPE !!!"
	}
}

// Returns neuron node type from its name
func NeuronTypeByName(name string) (NodeNeuronType, error) {
	switch name {
	case "HIDN":
		return HiddenNeuron, nil
	case "INPT":
		return InputNeuron, nil
	case "OUTP":
		return OutputNeuron, nil
	case "BIAS":
		return BiasNeuron, nil
	default:
		return math.MaxInt8, errors.New("Unknown neuron type name: " + name)
	}
}

// NodeActivationType defines the type of activation function to use for the neuron node
type NodeActivationType byte

// The neuron Activation function Types
const (
	// The sigmoid activation functions
	SigmoidPlainActivation NodeActivationType = iota + 1
	SigmoidReducedActivation
	SigmoidBipolarActivation
	SigmoidSteepenedActivation
	SigmoidApproximationActivation
	SigmoidSteepenedApproximationActivation
	SigmoidInverseAbsoluteActivation
	SigmoidLeftShiftedActivation
	SigmoidLeftShiftedSteepenedActivation
	SigmoidRightShiftedSteepenedActivation

	// The other activators assortment
	TanhActivation
	GaussianBipolarActivation
	LinearActivation
	LinearAbsActivation
	LinearClippedActivation
	NullActivation
	SignActivation
	SineActivation
	StepActivation

	// The modular activators (with multiple inputs/outputs)
	MultiplyModuleActivation
	MaxModuleActivation
	MinModuleActivation
)

// The neuron node activation function type
type ActivationFunction func(float64, []float64) float64
// The neurons module activation function type
type ModuleActivationFunction func([]float64, []float64) []float64

// The default node activators factory reference
var NodeActivators = NewNodeActivatorsFactory()

// The factory to provide appropriate neuron node activation function
type NodeActivatorsFactory struct {
	// The map of registered neuron node activators by type
	activators       map[NodeActivationType]ActivationFunction
	// The map of registered neurons module activators by type
	moduleActivators map[NodeActivationType]ModuleActivationFunction

	// The forward and inverse maps of activator type and function name
	forward          map[NodeActivationType]string
	inverse          map[string]NodeActivationType
}

// Returns node activator factory initialized with default activation functions
func NewNodeActivatorsFactory() *NodeActivatorsFactory {
	af := &NodeActivatorsFactory{
		activators:make(map[NodeActivationType]ActivationFunction),
		moduleActivators:make(map[NodeActivationType]ModuleActivationFunction),
		forward:make(map[NodeActivationType]string),
		inverse:make(map[string]NodeActivationType),
	}
	// Register neuron node activators
	af.Register(SigmoidPlainActivation, plainSigmoid, "SigmoidPlainActivation")
	af.Register(SigmoidReducedActivation, reducedSigmoid, "SigmoidReducedActivation")
	af.Register(SigmoidSteepenedActivation, steepenedSigmoid, "SigmoidSteepenedActivation")
	af.Register(SigmoidBipolarActivation, bipolarSigmoid, "SigmoidBipolarActivation")
	af.Register(SigmoidApproximationActivation, approximationSigmoid, "SigmoidApproximationActivation")
	af.Register(SigmoidSteepenedApproximationActivation, approximationSteepenedSigmoid, "SigmoidSteepenedApproximationActivation")
	af.Register(SigmoidInverseAbsoluteActivation, inverseAbsoluteSigmoid, "SigmoidInverseAbsoluteActivation")
	af.Register(SigmoidLeftShiftedActivation, leftShiftedSigmoid, "SigmoidLeftShiftedActivation")
	af.Register(SigmoidLeftShiftedSteepenedActivation, leftShiftedSteepenedSigmoid, "SigmoidLeftShiftedSteepenedActivation")
	af.Register(SigmoidRightShiftedSteepenedActivation, rightShiftedSteepenedSigmoid, "SigmoidRightShiftedSteepenedActivation")

	af.Register(TanhActivation, hyperbolicTangent, "TanhActivation")
	af.Register(GaussianBipolarActivation, bipolarGaussian, "GaussianBipolarActivation")
	af.Register(LinearActivation, linear, "LinearActivation")
	af.Register(LinearAbsActivation, absoluteLinear, "LinearAbsActivation")
	af.Register(LinearClippedActivation, clippedLinear, "LinearClippedActivation")
	af.Register(NullActivation, nullFunctor, "NullActivation")
	af.Register(SignActivation, signFunction, "SignActivation")
	af.Register(SineActivation, sineFunction, "SineActivation")
	af.Register(StepActivation, stepFunction, "StepActivation")

	// register neuron modules activators
	af.RegisterModule(MultiplyModuleActivation, multiplyModule, "MultiplyModuleActivation")
	af.RegisterModule(MaxModuleActivation, maxModule, "MaxModuleActivation")
	af.RegisterModule(MinModuleActivation, minModule, "MinModuleActivation")

	return af
}

// Method to calculate activation for specified neuron node based on it's ActivationType field value.
// Will return error and set -0.0 activation if unsupported activation type requested.
func (a *NodeActivatorsFactory) ActivateNode(node *NNode) (err error) {
	out, err := a.ActivateByType(node.ActivationSum, node.Params, node.ActivationType)
	if err == nil {
		node.setActivation(out)
	}
	return err
}

// Method to calculate activation value for give input and auxiliary parameters using activation function with specified type.
// Will return error and -0.0 activation if unsupported activation type requested.
func (a *NodeActivatorsFactory) ActivateByType(input float64, aux_params[]float64, a_type NodeActivationType) (float64, error) {
	if fn, ok := a.activators[a_type]; ok {
		return fn(input, aux_params), nil
	} else {
		return -0.0, errors.New(fmt.Sprintf("Unknown neuron activation type: %d", a_type))
	}
}

// Method to activate neuron module presented by provided node. As a result of execution the activation values of all
// input nodes will be processed by corresponding activation function and corresponding activation values of output nodes
// will be set. Will panic if unsupported activation type requested.
func (a *NodeActivatorsFactory) ActivateModule(module *NNode) error {
	inputs := make([]float64, len(module.Incoming))
	for i, v := range module.Incoming {
		inputs[i] = v.InNode.GetActiveOut()
	}

	outputs, err := a.ActivateModuleByType(inputs, module.Params, module.ActivationType)
	if err != nil {
		return err
	}
	if len(outputs) != len(module.Outgoing) {
		return errors.New(fmt.Sprintf(
			"The number of output parameters [%d] returned by module activator doesn't match " +
				"the number of output neurons of the module [%d]", len(outputs), len(module.Outgoing)))
	}
	// set outputs
	for i, out := range outputs {
		module.Outgoing[i].OutNode.setActivation(out)
		module.Outgoing[i].OutNode.isActive = true // activate output node
	}
	return nil
}

// Method will apply corresponding module activation function to the input values and returns appropriate output values.
// Will panic if unsupported activation function requested
func (a *NodeActivatorsFactory) ActivateModuleByType(inputs[] float64, aux_params[]float64, a_type NodeActivationType) ([]float64, error) {
	if fn, ok := a.moduleActivators[a_type]; ok {
		return fn(inputs, aux_params), nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown module activation type: %d", a_type))
	}
}

// Registers given neuron activation function with provided type and name into the factory
func (a *NodeActivatorsFactory) Register(a_type NodeActivationType, a_func ActivationFunction, f_name string) {
	// store function
	a.activators[a_type] = a_func
	// store name<->type bi-directional mapping
	a.forward[a_type] = f_name
	a.inverse[f_name] = a_type
}

// Registers given neuron module activation function with provided type and name into the factory
func (a *NodeActivatorsFactory) RegisterModule(a_type NodeActivationType, a_func ModuleActivationFunction, f_name string) {
	// store function
	a.moduleActivators[a_type] = a_func
	// store name<->type bi-directional mapping
	a.forward[a_type] = f_name
	a.inverse[f_name] = a_type
}

// Parse node activation type name and return corresponding activation type
func (a *NodeActivatorsFactory) ActivationTypeFromName(name string) (NodeActivationType, error) {
	if t, ok := a.inverse[name]; ok {
		return t, nil
	} else {
		return math.MaxInt8, errors.New("Unsupported activation type name: " + name)
	}
}

// Returns activation function name from given type
func (a *NodeActivatorsFactory) ActivationNameFromType(atype NodeActivationType) (string, error) {
	if n, ok := a.forward[atype]; ok {
		return n, nil
	} else {
		return "", errors.New(fmt.Sprintf("Unsupported activation type: %d", atype))
	}
}

// The sigmoid activation functions
var (
	// The plain sigmoid
	plainSigmoid = func(input float64, aux_params[]float64) float64 {
		return (1 / (1 + math.Exp(-input)))
	}
	// The reduced sigmoid
	reducedSigmoid = func(input float64, aux_params[]float64) float64 {
		return (1 / (1 + math.Exp(-0.5 * input)))
	}
	// The steepened sigmoid
	steepenedSigmoid = func(input float64, aux_params[]float64) float64 {
		return 1.0 / (1.0 + math.Exp(-4.924273 * input))
	}
	// The bipolar sigmoid activation function xrange->[-1,1] yrange->[-1,1]
	bipolarSigmoid = func(input float64, aux_params[]float64) float64 {
		return (2.0 / (1.0 + math.Exp(-4.924273 * input))) - 1.0
	}
	// The approximation sigmoid with squashing range [-4.0; 4.0]
	approximationSigmoid = func(input float64, aux_params[]float64) float64 {
		four, one_32nd := float64(4.0), float64(0.03125)
		if input < -4.0 {
			return 0.0
		} else if input < 0.0 {
			return (input + four) * (input + four) * one_32nd
		} else if input < 4.0 {
			return 1.0 - (input - four) * (input - four) * one_32nd
		} else {
			return 1.0
		}
	}
	// The steepened aproximation sigmoid with squashing range [-1.0; 1.0]
	approximationSteepenedSigmoid = func(input float64, aux_params[]float64) float64 {
		one, one_half := 1.0, 0.5
		if input < -1.0 {
			return 0.0
		} else if input < 0.0 {
			return (input + one) * (input + one) * one_half
		} else if input < 1.0 {
			return 1.0 - (input - one) * (input - one) * one_half
		} else {
			return 1.0;
		}
	}
	// The inverse absolute sigmoid
	inverseAbsoluteSigmoid = func(input float64, aux_params[]float64) float64 {
		return 0.5 + (input / (1.0 + math.Abs(input))) * 0.5
	}

	// The left/right shifted sigmoids
	leftShiftedSigmoid = func(input float64, aux_params[]float64) float64 {
		return 1.0 / (1.0 + math.Exp(-input - 2.4621365))
	}
	leftShiftedSteepenedSigmoid = func(input float64, aux_params[]float64) float64 {
		return 1.0 / (1.0 + math.Exp(-(4.924273 * input + 2.4621365)))
	}
	rightShiftedSteepenedSigmoid = func(input float64, aux_params[]float64) float64 {
		return 1.0 / (1.0 + math.Exp(-(4.924273 * input - 2.4621365)))
	}
)

// The other activation functions
var (
	// The hyperbolic tangent
	hyperbolicTangent = func(input float64, aux_params[]float64) float64 {
		return math.Tanh(0.9 * input)
	}
	// The bipolar Gaussian activator xrange->[-1,1] yrange->[-1,1]
	bipolarGaussian = func(input float64, aux_params[]float64) float64 {
		return 2.0 * math.Exp(-math.Pow(input * 2.5, 2.0)) - 1.0
	}
	// The absolute linear
	absoluteLinear = func(input float64, aux_params[]float64) float64 {
		return math.Abs(input)
	}
	// Linear activation function with clipping. By 'clipping' we mean the output value is linear between
	/// x = -1 and x = 1. Below -1 and above +1 the output is clipped at -1 and +1 respectively
	clippedLinear = func(input float64, aux_params[]float64) float64 {
		if (input < -1.0) {
			return -1.0
		}
		if (input > 1.0) {
			return 1.0
		}
		return input
	}
	// The linear activation
	linear = func(input float64, aux_params[]float64) float64 {
		return input
	}
	// The null activator
	nullFunctor = func(input float64, aux_params[]float64) float64 {
		return 0.0
	}
	// The sign activator
	signFunction = func(input float64, aux_params[]float64) float64 {
		if math.IsNaN(input) || input == 0.0 {
			return 0.0
		} else if math.Signbit(input) {
			return -1.0
		} else {
			return 1.0
		}
	}
	// The sine periodic activation with doubled period
	sineFunction = func(input float64, aux_params[]float64) float64 {
		return math.Sin(2.0 * input)
	}
	// The step function x<0 ? 0.0 : 1.0
	stepFunction = func(input float64, aux_params[]float64) float64 {
		if math.Signbit(input) {
			return 0.0
		} else {
			return 1.0
		}
	}
)

// The modular activators
var (
	// Multiplies input values and returns multiplication result
	multiplyModule = func(inputs []float64, aux_params[]float64) []float64 {
		ret := 1.0
		for _, v := range inputs {
			ret *= v
		}
		return []float64{ret}
	}
	// Finds maximal value among inputs and return it
	maxModule = func(inputs []float64, aux_params[]float64) []float64 {
		max := float64(math.MinInt64)
		for _, v := range inputs {
			max = math.Max(max, v)
		}
		return []float64{max}
	}
	// Finds minimal value among inputs and returns it
	minModule = func(inputs []float64, aux_params[]float64) []float64 {
		min := math.MaxFloat64
		for _, v := range inputs {
			min = math.Min(min, v)
		}
		return []float64{min}
	}
)
