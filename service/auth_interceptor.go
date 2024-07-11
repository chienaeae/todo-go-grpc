package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:      jwtManager,
		accessibleRoles: accessibleRoles,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (res any, err error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		claims, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, userClaimsKey, claims)

		return handler(ctx, req)
	}
}

type WrappedServerStream struct {
	grpc.ServerStream
	wrappedCtx context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.wrappedCtx
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		claims, err := interceptor.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		ctx := context.WithValue(ss.Context(), userClaimsKey, claims)
		return handler(srv, &WrappedServerStream{
			ServerStream: ss,
			wrappedCtx:   ctx,
		})
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (*UserClaims, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		// everyone can access the method
		return nil, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if claims.Role == role {
			return claims, nil
		}
	}

	return nil, status.Errorf(codes.Unauthenticated, "no permission to access this RPC")
}
