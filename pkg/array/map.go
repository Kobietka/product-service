package array

type TransformFunc[I, O any] func(I) O

func MapArray[I, O any](input []I, transform TransformFunc[I, O]) []O {
	output := make([]O, len(input))
	for index := range input {
		output[index] = transform(input[index])
	}
	return output
}
