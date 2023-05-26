package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"strings"
)

const (
	zeroFrame          = 0
	frameIndexIncrease = 2
	targetFrameIndex   = 3
	targetSplitNumber  = 2
	targetFrameNumber  = 1
)

func getCallerName(targetFrameIndex int) string {
	programCounters := make([]uintptr, targetFrameIndex+frameIndexIncrease)
	n := runtime.Callers(zeroFrame, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, zeroFrame; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}
	return frame.Function
}

func GetSpanName() string {
	fnName := getCallerName(targetFrameIndex)
	splitted := strings.Split(fnName, ".")
	if len(splitted) > targetSplitNumber {
		fnName = fmt.Sprintf("%s @ %s", strings.Join(splitted[(len(splitted)-targetSplitNumber):], "."),
			strings.Join(splitted[:(len(splitted)-targetSplitNumber)], "."))
	}

	return fnName
}

// GetMethodName projectGithubAddress = "github.com/inDriver/truck-api"
func GetMethodName(projectGithubAddress string) string {
	return strings.Replace(getCallerName(targetFrameIndex), projectGithubAddress, "", targetFrameNumber)
}

// StartSpan добавляет в существующую трассировку новую ветку
func StartSpan(parentCtx context.Context, spanName string) (ctx context.Context, span trace.Span, debugID string) {
	ctx, span = trace.SpanFromContext(parentCtx).TracerProvider().Tracer("").Start(parentCtx, spanName)
	traceID := span.SpanContext().TraceID().String()
	return ctx, span, traceID
}

// RegisterError регистрирует ошибку в трассировке
func RegisterError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
