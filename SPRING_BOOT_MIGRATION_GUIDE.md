# Greenlight: From Go to Idiomatic Spring Boot
## A Production-Grade Project-Based Learning Guide

Welcome to the **Greenlight Spring Boot Migration Guide**. This interactive, hands-on "book" translates your Go-based Greenlight API (from Alex Edwards' *Let's Go Further*) into a modern, idiomatic, production-grade **Spring Boot 3 (Java 21)** backend application.

Since you have already read **Spring Starts Here** through Chapter 7, you understand:
- **Inversion of Control (IoC)** and the **Spring ApplicationContext**
- **Beans** and the `@Bean` / `@Component` annotations
- **Dependency Injection (DI)** with `@Autowired` and constructor injection
- **Bean Wiring & Scopes**

Now, we bridge the gap between Spring core concepts and building real-world RESTful services using **Spring Boot**, **Spring Data JPA**, **Spring Security**, and **Spring Validation**.

---

## How to Use This Guide
In each step of each module, you will find:
1. **Target File Path:** The exact path relative to your project root (e.g. `src/main/java/com/greenlight/api/system/HealthController.java`).
2. **Architectural Rationale:** *Why* Spring does it this way compared to Go.
3. **Template Blueprint:** A code skeleton with the exact package, annotations, constructor injection, and method signatures so you can write the implementation without guessing where files go.
4. **Verification & Git Milestone:** How to test and commit your progress.

> [!TIP]
> **Where the Code Starts:** This document is an extensive 2,550+ line comprehensive guide. Lines 1–157 cover the architectural paradigm shift, module roadmap, and Module 0. **The step-by-step code snippets, file paths, and class blueprints begin at Module 1 on line 158.**

---

## Table of Contents
1. [The Architectural Paradigm Shift (Go vs. Spring Boot)](#chapter-1-the-architectural-paradigm-shift-go-vs-spring-boot)
2. [Module Roadmap & Git Strategy](#chapter-2-module-roadmap--git-strategy)
3. [Module 0: Project Initialization & IntelliJ Setup](#module-0-project-initialization--intellij-setup)
4. [Module 1: Foundations, Configuration & System Health](#module-1-foundations-configuration--system-health)
5. [Module 2: Domain Modeling, DTOs & Jakarta Validation](#module-2-domain-modeling-dtos--jakarta-validation)
6. [Module 3: Global Error Handling & Problem Details](#module-3-global-error-handling--problem-details)
7. [Module 4: Persistence with PostgreSQL, Flyway & Spring Data JPA](#module-4-persistence-with-postgresql-flyway--spring-data-jpa)
8. [Module 5: Optimistic Locking & Concurrency Control](#module-5-optimistic-locking--concurrency-control)
9. [Module 6: Advanced Querying — Pagination, Dynamic Sorting & Search](#module-6-advanced-querying--pagination-dynamic-sorting--search)
10. [Module 7: User Management & Password Security](#module-7-user-management--password-security)
11. [Module 8: Asynchronous Processing & Email Dispatch](#module-8-asynchronous-processing--email-dispatch)
12. [Module 9: Authentication & Stateful Bearer Tokens via Spring Security Filter Chain](#module-9-authentication--stateful-bearer-tokens-via-spring-security-filter-chain)
13. [Module 10: Role-Based Access Control (RBAC) & Method Security](#module-10-role-based-access-control-rbac--method-security)
14. [Module 11: Rate Limiting & CORS Policies](#module-11-rate-limiting--cors-policies)
15. [Module 12: Production Observability, Metrics & Actuator](#module-12-production-observability-metrics--actuator)
16. [Masterclass: Comprehensive Testing Strategy (Unit, Slice & Integration)](#chapter-16-masterclass-comprehensive-testing-strategy)
17. [Appendix: Go to Spring Boot Rosetta Stone](#appendix-go-to-spring-boot-rosetta-stone)

---

## Chapter 1: The Architectural Paradigm Shift (Go vs. Spring Boot)

In Go (`greenlight`), you built an explicit, procedural server:
```
main.go (creates db, config, logger)
   │
   ▼
app := &application{config, logger, models, mailer, wg} (explicit dependency bag)
   │
   ▼
routes.go (manual httprouter registration + middleware wrapper onion)
   │
   ▼
handlers -> models (manual SQL scanning) -> JSON helpers
```

In an idiomatic **Spring Boot** application:
1. **The ApplicationContext is the Container:** You never create a monolithic `application` struct. Classes declare their dependencies via constructor injection, and the Spring IoC container wires them automatically.
2. **Layered Architecture (Separation of Concerns):**
   - **Controller Layer (`@RestController`):** Handles HTTP mapping, JSON serialization/deserialization, query param parsing, and delegates to the Service.
   - **Service Layer (`@Service`):** Houses business logic, transaction boundaries (`@Transactional`), domain validation, and coordinates repositories and mailers.
   - **Repository Layer (`@Repository`):** Encapsulates data access using Spring Data JPA abstractions or custom JDBC templates.
   - **Data Transfer Objects (DTOs):** In Go, you often attached JSON tags directly to database structs or used anonymous structs. In production Spring, you **never** expose JPA database entities directly to HTTP clients. You use Java 21 `record`s as Request/Response DTOs.
3. **Cross-Cutting Concerns via Interceptors / Filters / AOP:**
   - In Go, you manually nested middleware closures: `app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))`.
   - In Spring, cross-cutting concerns use standard Filter Chains (Spring Security), `@RestControllerAdvice` (exception recovery), and Spring WebMvc configurers (CORS).

---

## Chapter 2: Module Roadmap & Git Strategy

Each module represents a self-contained learning unit. To develop the muscle memory of a professional backend engineer:
1. **Create a feature branch** for the module (e.g., `git checkout -b feat/module-1-foundations`).
2. **Implement the classes** in IntelliJ IDEA.
3. **Verify** using `curl` or IntelliJ's built-in HTTP Client (`.http` files).
4. **Commit and merge** into `main`, then push to your remote repository.

| Module | Topic | Go greenlight Equivalent | Spring Boot Concepts |
| :--- | :--- | :--- | :--- |
| **0** | Project Setup | `go.mod`, `Makefile` | Spring Initializr, Maven, Java 21, IntelliJ |
| **1** | Foundations & Health | `main.go`, `healthcheck.go` | `@SpringBootApplication`, `@RestController`, ConfigurationProperties |
| **2** | Models & Validation | `internal/validator`, `movies.go` DTOs | Jakarta Bean Validation (`@Valid`), Java 21 Records |
| **3** | Error Handling | `errors.go` | `@RestControllerAdvice`, `@ExceptionHandler`, unified envelopes |
| **4** | Persistence & Flyway | `migrations/`, `internal/data` | Flyway, Spring Data JPA, HikariCP, PostgreSQL |
| **5** | Optimistic Locking | `version = version + 1`, `409 Conflict` | `@Version`, `OptimisticLockingFailureException` |
| **6** | Querying & Pagination | `GetAll(title, genres, filters)` | `Pageable`, `Specification` / dynamic queries, Page metadata |
| **7** | Users & Passwords | `users.go`, bcrypt | `BCryptPasswordEncoder`, Unique constraints, User Service |
| **8** | Async & Mailer | `mailer.go`, `wg.Add(1)` background | `@EnableAsync`, `@Async`, `ThreadPoolTaskExecutor`, Spring Events |
| **9** | Auth & Token Filter | `authenticate`, `tokens.go` | Spring Security 6 `SecurityFilterChain`, `OncePerRequestFilter` |
| **10** | RBAC & Permissions | `requirePermission()`, `users_permissions` | `@PreAuthorize`, `GrantedAuthority`, Spring Method Security |
| **11** | Rate Limiter & CORS | `rateLimit()`, `enableCORS()` | Bucket4j / Filter, `CorsConfigurationSource` |
| **12** | Observability | `expvar`, `/debug/vars` | Spring Boot Actuator, Micrometer Metrics |
| **Final** | Testing Masterclass | `httptest` (Go Chapter 20) | JUnit 5, Mockito, `@WebMvcTest`, `@DataJpaTest`, `@SpringBootTest` |

---

## Module 0: Project Initialization & IntelliJ Setup

### Objective
Generate the project skeleton, set up IntelliJ IDEA, understand the directory structure, and initialize Git.

### 1. Generating the Project with Spring Initializr
Open [start.spring.io](https://start.spring.io) (or use IntelliJ's built-in **File -> New -> Project -> Spring Initializr**):
- **Project:** Maven
- **Language:** Java
- **Spring Boot:** `3.3.x` (or latest stable 3.x)
- **Group:** `com.greenlight`
- **Artifact:** `greenlight-api`
- **Name:** `greenlight-api`
- **Package name:** `com.greenlight.api`
- **Packaging:** Jar
- **Java version:** 21

### Dependencies to Select in Initializr:
- **Spring Web** (Embedded Tomcat, Spring MVC, REST)
- **Spring Data JPA** (Hibernate, EntityManager, Spring Repositories)
- **PostgreSQL Driver** (JDBC driver for PostgreSQL)
- **Flyway Migration** (Database migration versioning)
- **Validation** (Jakarta Bean Validation / Hibernate Validator)
- **Spring Security** (Authentication, Authorization, Filter chains)
- **Spring Boot Actuator** (Production metrics, health check)

#### Additional Dependencies to Add in `pom.xml`:
Open `pom.xml` in IntelliJ and add these two essential libraries inside `<dependencies>`:

```xml
<!-- 1. Auto-load .env file into Spring Environment (replaces Go's godotenv.Load()) -->
<dependency>
    <groupId>me.paulschwarz</groupId>
    <artifactId>spring-dotenv</artifactId>
    <version>4.0.0</version>
</dependency>

<!-- 2. Official Resend Java SDK for sending emails -->
<dependency>
    <groupId>com.resend</groupId>
    <artifactId>resend-java</artifactId>
    <version>3.1.0</version>
</dependency>
```

### 2. Package Structure Conventions
In IntelliJ, structure your `src/main/java/com/greenlight/api/` as follows:
```text
src/main/java/com/greenlight/api/
├── GreenlightApplication.java        // Main entry point
├── config/                           // AppProperties, DatabaseConfig, MailConfig, SecurityConfig, AsyncConfig
├── common/                           // Shared utilities, envelopes, custom exceptions
│   ├── dto/
│   ├── error/
│   └── validator/
├── system/                           // System healthcheck & metadata
│   ├── HealthController.java
│   └── dto/
│       └── HealthCheckResponse.java
├── movie/                            // Movie feature module (Entity, Repo, Service, Controller, DTOs)
│   └── dto/
├── user/                             // User feature module
│   ├── dto/
│   └── event/
├── token/                            // Token authentication module
│   └── dto/
└── mail/                             // Email notification module
    └── MailService.java
```

### 3. Git Milestone
```bash
git init
git add .
git commit -m "chore: initialize spring boot 3 project with maven, java 21, and dependencies"
```

---

## Module 1: Foundations, Neon PostgreSQL Connection, Resend Mailer & System Health

### What you are migrating from Go:
- `cmd/api/main.go`:
  - Loading `.env` via `godotenv.Load()`.
  - Establishing PostgreSQL connection pool (`openDB()`) with Neon.
  - Configuring Resend client (`mailer.New()`).
  - Graceful connection ping (`db.PingContext()`).
- `cmd/api/healthcheck.go`: `/v1/healthcheck` endpoint returning system status and DB health.

### Spring Boot Concepts to Master:
1. **`@SpringBootApplication`**: Bootstraps the ApplicationContext and component scanning.
2. **Neon Database & HikariCP**: In Go, you configured `maxOpenConns: 25`, `maxIdleConns: 12`, `maxIdleTime: 15m`. In Spring Boot, HikariCP is the high-performance connection pool managed as a `DataSource` bean.
3. **Resend Java SDK Integration**: Creating a `@Bean` for the `Resend` client reading `RESEND_API_KEY`.
4. **Auto-Loading `.env`**: Using `spring-dotenv` so your existing `.env` credentials work seamlessly.
5. **Testing DB Connectivity in HealthCheck**: Injecting `DataSource` to verify the Neon connection pool is active.

---

### Step-by-Step Implementation Guide

#### Step 1: Ensure `.env` is in the Project Root
- **File:** `.env` (in your Spring Boot project root directory)
- **Purpose:** Environment variables matching your existing Go setup.

```env
DATABASE_URL=''
RESEND_API_KEY=
```

> [!NOTE]
> The `spring-dotenv` dependency automatically loads these variables into Spring's environment upon startup!

---

#### Step 2: Configure `application.yml`
- **File:** `src/main/resources/application.yml`
- **Purpose:** Centralized application configuration.

```yaml
server:
  port: 4000
  shutdown: graceful

spring:
  application:
    name: greenlight-api
  profiles:
    active: dev
  lifecycle:
    timeout-per-shutdown-phase: 20s
  jpa:
    open-in-view: false
    hibernate:
      ddl-auto: none # Disabled for now; Flyway will manage schema in Module 4
    show-sql: false
  flyway:
    enabled: false # Temporarily disabled until we add migrations in Module 4

# Custom Greenlight Application Properties
app:
  version: "1.0.0"
  env: "development"
  mailer:
    sender: "Greenlight <no-reply@sohail.world>"
```

---

#### Step 3: Implement `DatabaseConfig` for Neon PostgreSQL & HikariCP
- **File:** `src/main/java/com/greenlight/api/config/DatabaseConfig.java`
- **Purpose:** Configures the `HikariDataSource` connection pool. Handles Neon connection URLs (converting standard URI `postgresql://` to JDBC `jdbc:postgresql://`) and sets the exact connection pool limits from Go.

```java
package com.greenlight.api.config;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.sql.DataSource;
import java.net.URI;

@Configuration
public class DatabaseConfig {

    @Value("${DATABASE_URL:}")
    private String databaseUrl;

    @Bean
    public DataSource dataSource() {
        HikariConfig config = new HikariConfig();

        // Neon connection string parsing
        if (databaseUrl.startsWith("postgresql://") || databaseUrl.startsWith("postgres://")) {
            // Strip out channel_binding parameter if present as pgjdbc handles sslmode directly
            String cleanUrl = databaseUrl.replace("&channel_binding=require", "")
                                         .replace("?channel_binding=require&", "?")
                                         .replace("?channel_binding=require", "");

            URI uri = URI.create(cleanUrl);
            String host = uri.getHost();
            int port = uri.getPort() == -1 ? 5432 : uri.getPort();
            String path = uri.getPath();
            String query = uri.getQuery() != null ? "?" + uri.getQuery() : "?sslmode=require";

            config.setJdbcUrl("jdbc:postgresql://" + host + ":" + port + path + query);

            if (uri.getUserInfo() != null) {
                String[] credentials = uri.getUserInfo().split(":");
                config.setUsername(credentials[0]);
                if (credentials.length > 1) {
                    config.setPassword(credentials[1]);
                }
            }
        } else if (databaseUrl.startsWith("jdbc:")) {
            config.setJdbcUrl(databaseUrl);
        } else {
            throw new IllegalStateException("DATABASE_URL is not configured or invalid!");
        }

        // Match Go connection pool settings exactly:
        // maxOpenConns: 25, maxIdleConns: 12, maxIdleTime: 15m, timeout: 5s
        config.setMaximumPoolSize(25);
        config.setMinimumIdle(12);
        config.setIdleTimeout(900_000); // 15 minutes in ms
        config.setConnectionTimeout(5_000); // 5 seconds
        config.setPoolName("GreenlightHikariPool");

        return new HikariDataSource(config);
    }
}
```

---

#### Step 4: Configure Resend Mailer Bean & `MailService`
- **File:** `src/main/java/com/greenlight/api/config/MailConfig.java`
- **Purpose:** Initializes the Resend API client using `RESEND_API_KEY`.

```java
package com.greenlight.api.config;

import com.resend.Resend;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class MailConfig {

    @Bean
    public Resend resend(@Value("${RESEND_API_KEY:}") String apiKey) {
        if (apiKey == null || apiKey.isBlank()) {
            throw new IllegalStateException("RESEND_API_KEY must be provided in .env");
        }
        return new Resend(apiKey);
    }
}
```

- **File:** `src/main/java/com/greenlight/api/mail/MailService.java`
- **Purpose:** Service for sending emails via Resend.

```java
package com.greenlight.api.mail;

import com.greenlight.api.config.AppProperties;
import com.resend.Resend;
import com.resend.core.exception.ResendException;
import com.resend.services.emails.model.CreateEmailOptions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

@Service
public class MailService {

    private static final Logger log = LoggerFactory.getLogger(MailService.class);

    private final Resend resend;
    private final AppProperties appProperties;

    public MailService(Resend resend, AppProperties appProperties) {
        this.resend = resend;
        this.appProperties = appProperties;
    }

    public void sendWelcomeEmail(String recipientEmail, String recipientName, String activationToken) {
        try {
            CreateEmailOptions params = CreateEmailOptions.builder()
                .from(appProperties.mailer().sender())
                .to(recipientEmail)
                .subject("Welcome to Greenlight!")
                .html("<p>Hi " + recipientName + ",</p><p>Please activate your account with token: <strong>" + activationToken + "</strong></p>")
                .build();

            resend.emails().send(params);
            log.info("Welcome email sent to {}", recipientEmail);
        } catch (ResendException e) {
            log.error("Failed to send welcome email via Resend: ", e);
        }
    }
}
```

---

#### Step 5: Create Type-Safe `AppProperties` Record
- **File:** `src/main/java/com/greenlight/api/config/AppProperties.java`
- **Purpose:** Binds application properties including mailer settings.

```java
package com.greenlight.api.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app")
public record AppProperties(
    String version,
    String env,
    MailerProperties mailer
) {
    public record MailerProperties(String sender) {}
}
```

---

#### Step 6: Enable Configuration Properties Scan
- **File:** `src/main/java/com/greenlight/api/GreenlightApplication.java`

```java
package com.greenlight.api;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;

@SpringBootApplication
@ConfigurationPropertiesScan
public class GreenlightApplication {

    public static void main(String[] args) {
        SpringApplication.run(GreenlightApplication.class, args);
    }
}
```

---

#### Step 7: Create HealthCheck Response DTO
- **File:** `src/main/java/com/greenlight/api/system/dto/HealthCheckResponse.java`
- **Purpose:** Response payload including database connection status.

```java
package com.greenlight.api.system.dto;

public record HealthCheckResponse(
    String status,
    String environment,
    String version,
    String database,
    String mailer
) {}
```

---

#### Step 8: Implement `HealthController` with Database Ping
- **File:** `src/main/java/com/greenlight/api/system/HealthController.java`
- **Purpose:** Pings Neon database connection pool on demand (replaces Go's `db.PingContext()`).

```java
package com.greenlight.api.system;

import com.greenlight.api.config.AppProperties;
import com.greenlight.api.system.dto.HealthCheckResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.SQLException;

@RestController
@RequestMapping("/v1")
public class HealthController {

    private final AppProperties appProperties;
    private final DataSource dataSource;

    public HealthController(AppProperties appProperties, DataSource dataSource) {
        this.appProperties = appProperties;
        this.dataSource = dataSource;
    }

    @GetMapping("/healthcheck")
    public ResponseEntity<HealthCheckResponse> healthcheck() {
        String dbStatus = "connected";

        // Ping database connection with 3-second timeout (identical to Go db.PingContext)
        try (Connection conn = dataSource.getConnection()) {
            if (!conn.isValid(3)) {
                dbStatus = "unavailable";
            }
        } catch (SQLException e) {
            dbStatus = "error: " + e.getMessage();
        }

        HealthCheckResponse response = new HealthCheckResponse(
            "available",
            appProperties.env(),
            appProperties.version(),
            dbStatus,
            "configured"
        );

        return ResponseEntity.ok(response);
    }
}
```

---

### Verification
1. Ensure your `.env` file is in your project root with your Neon `DATABASE_URL` and `RESEND_API_KEY`.
2. Start `GreenlightApplication` in IntelliJ IDEA.
   Notice the console log:
   `HikariPool-1 - Added connection org.postgresql.jdbc.PgConnection@...`
   `GreenlightHikariPool - Start completed.`
3. Execute in terminal:
```bash
curl -i http://localhost:4000/v1/healthcheck
```
Expected output:
```json
{
  "status": "available",
  "environment": "development",
  "version": "1.0.0",
  "database": "connected",
  "mailer": "configured"
}
```

### Git Milestone
```bash
git checkout -b feat/module-1-foundations
git add .
git commit -m "feat(system): configure neon postgresql hikari pool, resend mailer, and healthcheck"
git checkout main && git merge feat/module-1-foundations
git push origin main
```

---

## Module 2: Domain Modeling, DTOs & Jakarta Validation

### What you are migrating from Go:
- `internal/validator/validator.go`: Custom map-based validator (`v.Check(movie.Year >= 1888, ...)`).
- `internal/data/movies.go`: `ValidateMovie(v *validator.Validator, movie *Movie)`:
  - Title not blank, max 500 bytes.
  - Year >= 1888 and not in the future.
  - Runtime positive integer.
  - Genres: 1 to 5 unique items.
  - Image URL: valid URL format.
  - Synopsis: 10 to 1000 characters.

### Spring Boot Concepts to Master:
1. **Jakarta Bean Validation (JSR 380):** Declarative validation using annotations (`@NotBlank`, `@Min`, `@Max`, `@Size`, `@Positive`).
2. **Java 21 Records for DTOs:** Immutable, clean data carriers without getter/setter boilerplate.
3. **Custom Constraint Validator:** Creating annotations like `@UniqueElements` to validate custom constraints.
4. **Triggering Validation:** The `@Valid` annotation in `@RequestBody`.

---

### Step-by-Step Implementation Guide

#### Step 1: Create Custom `@UniqueElements` Validation Annotation
- **File:** `src/main/java/com/greenlight/api/common/validator/UniqueElements.java`
- **Purpose:** Custom constraint annotation for collections.

```java
package com.greenlight.api.common.validator;

import jakarta.validation.Constraint;
import jakarta.validation.Payload;
import java.lang.annotation.*;

@Documented
@Constraint(validatedBy = UniqueElementsValidator.class)
@Target({ElementType.FIELD, ElementType.PARAMETER})
@Retention(RetentionPolicy.RUNTIME)
public @interface UniqueElements {
    String message() default "must not contain duplicate values";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}
```

---

#### Step 2: Implement `UniqueElementsValidator`
- **File:** `src/main/java/com/greenlight/api/common/validator/UniqueElementsValidator.java`
- **Purpose:** Validator logic checking for duplicates using a `Set`.

```java
package com.greenlight.api.common.validator;

import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;
import java.util.Collection;
import java.util.HashSet;

public class UniqueElementsValidator implements ConstraintValidator<UniqueElements, Collection<?>> {

    @Override
    public boolean isValid(Collection<?> collection, ConstraintValidatorContext context) {
        if (collection == null || collection.isEmpty()) {
            return true; // Use @NotEmpty or @NotNull to validate presence
        }
        // TODO: Return true if all elements are distinct
        return new HashSet<>(collection).size() == collection.size();
    }
}
```

---

#### Step 3: Create `CreateMovieRequest` DTO
- **File:** `src/main/java/com/greenlight/api/movie/dto/CreateMovieRequest.java`
- **Purpose:** Strongly-typed, validated incoming payload for creating a movie.

```java
package com.greenlight.api.movie.dto;

import com.greenlight.api.common.validator.UniqueElements;
import jakarta.validation.constraints.*;
import org.hibernate.validator.constraints.URL;
import java.util.List;

public record CreateMovieRequest(
    @NotBlank(message = "must be provided")
    @Size(max = 500, message = "must not be more than 500 bytes long")
    String title,

    @NotNull(message = "must be provided")
    @Min(value = 1888, message = "must be greater than 1888")
    Integer year,

    @NotNull(message = "must be provided")
    @Positive(message = "must be a positive integer")
    Integer runtime,

    @NotNull(message = "must be provided")
    @Size(min = 1, max = 5, message = "must contain between 1 and 5 genres")
    @UniqueElements(message = "must not contain duplicate values")
    List<String> genres,

    @NotBlank(message = "must be provided")
    @URL(message = "must be a valid URL")
    String imageUrl,

    @NotBlank(message = "must be provided")
    @Size(min = 10, max = 1000, message = "must be between 10 and 1000 characters")
    String synopsis
) {}
```

---

#### Step 4: Create Skeleton `MovieController`
- **File:** `src/main/java/com/greenlight/api/movie/MovieController.java`
- **Purpose:** REST controller exposing `/v1/movies` endpoints.

```java
package com.greenlight.api.movie;

import com.greenlight.api.movie.dto.CreateMovieRequest;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/movies")
public class MovieController {

    @PostMapping
    public ResponseEntity<CreateMovieRequest> createMovie(@Valid @RequestBody CreateMovieRequest request) {
        // Temporary: echo back the validated request until persistence is wired in Module 4
        return ResponseEntity.status(HttpStatus.CREATED).body(request);
    }
}
```

### Git Milestone
```bash
git checkout -b feat/module-2-validation
git add .
git commit -m "feat(movie): implement create movie request DTO with jakarta validation"
git checkout main && git merge feat/module-2-validation
git push origin main
```

---

## Module 3: Global Error Handling & Problem Details

### What you are migrating from Go:
- `cmd/api/errors.go`:
  - `errorResponse()`: JSON envelope `{"error": ...}`.
  - `failedValidationResponse()`: returns 422 with map of `{field: error}`.
  - `notFoundResponse()`: returns 404.
  - `serverErrorResponse()`: returns 500.
  - `methodNotAllowedResponse()`: returns 405.

### Spring Boot Concepts to Master:
1. **`@RestControllerAdvice`**: Global interceptor for all exceptions thrown by any `@RestController`.
2. **`@ExceptionHandler`**: Catch specific exceptions (e.g., `MethodArgumentNotValidException`, `ResourceNotFoundException`).
3. **Consistent Error Envelope**: Matching Go's `{"error": ...}` JSON structure.

---

### Step-by-Step Implementation Guide

#### Step 1: Define the `ErrorResponse` Envelope Record
- **File:** `src/main/java/com/greenlight/api/common/error/ErrorResponse.java`
- **Purpose:** Standard JSON wrapper for error payloads.

```java
package com.greenlight.api.common.error;

public record ErrorResponse(Object error) {}
```

---

#### Step 2: Create Custom Domain Exceptions
- **File:** `src/main/java/com/greenlight/api/common/error/ResourceNotFoundException.java`
- **Purpose:** Thrown when an entity with a specific ID does not exist in the database.

```java
package com.greenlight.api.common.error;

public class ResourceNotFoundException extends RuntimeException {
    public ResourceNotFoundException(String message) {
        super(message);
    }
}
```

---

#### Step 3: Implement `GlobalExceptionHandler`
- **File:** `src/main/java/com/greenlight/api/common/error/GlobalExceptionHandler.java`
- **Purpose:** Captures exceptions, logs details, and formats standard HTTP responses.

```java
package com.greenlight.api.common.error;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.FieldError;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.servlet.resource.NoResourceFoundException;

import java.util.HashMap;
import java.util.Map;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    // 1. Validation errors (HTTP 422 Unprocessable Entity, matching Go's failedValidationResponse)
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorResponse> handleValidationExceptions(MethodArgumentNotValidException ex) {
        Map<String, String> errors = new HashMap<>();
        for (FieldError fieldError : ex.getBindingResult().getFieldErrors()) {
            errors.put(fieldError.getField(), fieldError.getDefaultMessage());
        }
        return ResponseEntity.status(HttpStatus.UNPROCESSABLE_ENTITY).body(new ErrorResponse(errors));
    }

    // 2. Resource Not Found (HTTP 404 Not Found)
    @ExceptionHandler({ResourceNotFoundException.class, NoResourceFoundException.class})
    public ResponseEntity<ErrorResponse> handleNotFound(Exception ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
            .body(new ErrorResponse("the requested resource could not be found"));
    }

    // 3. Method Not Supported (HTTP 405 Method Not Allowed)
    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    public ResponseEntity<ErrorResponse> handleMethodNotSupported(HttpRequestMethodNotSupportedException ex) {
        String message = String.format("the %s method is not supported for this resource", ex.getMethod());
        return ResponseEntity.status(HttpStatus.METHOD_NOT_ALLOWED).body(new ErrorResponse(message));
    }

    // 4. Catch-all for uncaught exceptions (HTTP 500 Internal Server Error)
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorResponse> handleGeneralException(Exception ex) {
        log.error("Internal server error: ", ex);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
            .body(new ErrorResponse("the server encountered a problem and could not process your request"));
    }
}
```

### Verification
Run the application and send an invalid payload:
```bash
curl -i -X POST http://localhost:4000/v1/movies \
  -H "Content-Type: application/json" \
  -d '{"title": "", "year": 1800, "genres": ["drama", "drama"]}'
```
Expected HTTP 422:
```json
{
  "error": {
    "title": "must be provided",
    "year": "must be greater than 1888",
    "genres": "must not contain duplicate values",
    "imageUrl": "must be provided",
    "synopsis": "must be provided",
    "runtime": "must be provided"
  }
}
```

### Git Milestone
```bash
git checkout -b feat/module-3-error-handling
git add .
git commit -m "feat(error): implement global exception handler with 422 validation response"
git checkout main && git merge feat/module-3-error-handling
git push origin main
```

---

## Module 4: Persistence with PostgreSQL, Flyway & Spring Data JPA

### What you are migrating from Go:
- `migrations/`: The raw SQL migration scripts (`000001_create_movies_table.up.sql`, etc.).
- `cmd/api/main.go` `openDB()`: PostgreSQL connection pool configuration.
- `internal/data/movies.go`: `Insert()`, `Get()`, `Delete()`.

### Spring Boot Concepts to Master:
1. **Flyway Migrations:** Versioned migration scripts placed in `src/main/resources/db/migration/` automatically executed on startup.
2. **JPA Entity (`@Entity`, `@Table`, `@Id`, `@GeneratedValue`):** Mapping Java classes to relational tables.
3. **Handling Collections:** Mapping `genres` using `@ElementCollection`.
4. **Spring Data JPA Repositories:** Interfaces extending `JpaRepository<Movie, Long>`.
5. **Declarative Transactions (`@Transactional`):** Automatic transaction management.

---

### Step-by-Step Implementation Guide

#### Step 1: Configure Datasource and Flyway in `application.yml`
- **File:** `src/main/resources/application.yml`

```yaml
spring:
  datasource:
    url: ${DATABASE_URL:jdbc:postgresql://localhost:5432/greenlight}
    username: ${DB_USER:greenlight}
    password: ${DB_PASS:pa55word}
    hikari:
      maximum-pool-size: 25
      minimum-idle: 12
      idle-timeout: 900000 # 15 minutes
      connection-timeout: 5000 # 5 seconds
  jpa:
    hibernate:
      ddl-auto: validate # Flyway manages schema, Hibernate strictly validates
    show-sql: false
    properties:
      hibernate.format_sql: true
  flyway:
    enabled: true
    locations: classpath:db/migration
```

---

#### Step 2: Port Migrations to Flyway
Place your SQL scripts in `src/main/resources/db/migration/`:
- **File:** `src/main/resources/db/migration/V1__create_movies_table.sql`
```sql
CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    image_url text NOT NULL,
    synopsis text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
```

- **File:** `src/main/resources/db/migration/V2__add_movies_check_constraints.sql`
```sql
ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year >= 1888);
ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
```

- **File:** `src/main/resources/db/migration/V3__add_movies_indexes.sql`
```sql
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);
```

---

#### Step 3: Implement `Movie` JPA Entity
- **File:** `src/main/java/com/greenlight/api/movie/Movie.java`
- **Purpose:** Relational database entity.

```java
package com.greenlight.api.movie;

import jakarta.persistence.*;
import org.hibernate.annotations.CreationTimestamp;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "movies")
public class Movie {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(nullable = false, length = 500)
    private String title;

    @Column(nullable = false)
    private Integer year;

    @Column(nullable = false)
    private Integer runtime;

    // PostgreSQL text[] mapping via collection table or converter
    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "movie_genres", joinColumns = @JoinColumn(name = "movie_id"))
    @Column(name = "genre")
    private List<String> genres = new ArrayList<>();

    @Column(name = "image_url", nullable = false)
    private String imageUrl;

    @Column(nullable = false, length = 1000)
    private String synopsis;

    @Version
    @Column(nullable = false)
    private Integer version = 1;

    // Constructors
    public Movie() {}

    // Getters and Setters
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public Instant getCreatedAt() { return createdAt; }
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    public Integer getYear() { return year; }
    public void setYear(Integer year) { this.year = year; }
    public Integer getRuntime() { return runtime; }
    public void setRuntime(Integer runtime) { this.runtime = runtime; }
    public List<String> getGenres() { return genres; }
    public void setGenres(List<String> genres) { this.genres = genres; }
    public String getImageUrl() { return imageUrl; }
    public void setImageUrl(String imageUrl) { this.imageUrl = imageUrl; }
    public String getSynopsis() { return synopsis; }
    public void setSynopsis(String synopsis) { this.synopsis = synopsis; }
    public Integer getVersion() { return version; }
    public void setVersion(Integer version) { this.version = version; }
}
```

---

#### Step 4: Create `MovieResponse` DTO
- **File:** `src/main/java/com/greenlight/api/movie/dto/MovieResponse.java`
- **Purpose:** Read-only representation returned to clients.

```java
package com.greenlight.api.movie.dto;

import com.greenlight.api.movie.Movie;
import java.util.List;

public record MovieResponse(
    Long id,
    String title,
    Integer year,
    Integer runtime,
    List<String> genres,
    String imageUrl,
    String synopsis,
    Integer version
) {
    public static MovieResponse fromEntity(Movie movie) {
        return new MovieResponse(
            movie.getId(),
            movie.getTitle(),
            movie.getYear(),
            movie.getRuntime(),
            movie.getGenres(),
            movie.getImageUrl(),
            movie.getSynopsis(),
            movie.getVersion()
        );
    }
}
```

---

#### Step 5: Implement `MovieRepository`
- **File:** `src/main/java/com/greenlight/api/movie/MovieRepository.java`
- **Purpose:** Spring Data repository interface.

```java
package com.greenlight.api.movie;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.stereotype.Repository;

@Repository
public interface MovieRepository extends JpaRepository<Movie, Long>, JpaSpecificationExecutor<Movie> {
}
```

---

#### Step 6: Implement `MovieService`
- **File:** `src/main/java/com/greenlight/api/movie/MovieService.java`
- **Purpose:** Handles transactional business logic.

```java
package com.greenlight.api.movie;

import com.greenlight.api.common.error.ResourceNotFoundException;
import com.greenlight.api.movie.dto.CreateMovieRequest;
import com.greenlight.api.movie.dto.MovieResponse;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@Transactional(readOnly = true)
public class MovieService {

    private final MovieRepository movieRepository;

    public MovieService(MovieRepository movieRepository) {
        this.movieRepository = movieRepository;
    }

    @Transactional
    public MovieResponse createMovie(CreateMovieRequest request) {
        Movie movie = new Movie();
        movie.setTitle(request.title());
        movie.setYear(request.year());
        movie.setRuntime(request.runtime());
        movie.setGenres(request.genres());
        movie.setImageUrl(request.imageUrl());
        movie.setSynopsis(request.synopsis());

        Movie saved = movieRepository.save(movie);
        return MovieResponse.fromEntity(saved);
    }

    public MovieResponse getMovieById(Long id) {
        return movieRepository.findById(id)
            .map(MovieResponse::fromEntity)
            .orElseThrow(() -> new ResourceNotFoundException("Movie not found"));
    }

    @Transactional
    public void deleteMovie(Long id) {
        if (!movieRepository.existsById(id)) {
            throw new ResourceNotFoundException("Movie not found");
        }
        movieRepository.deleteById(id);
    }
}
```

---

#### Step 7: Update `MovieController`
- **File:** `src/main/java/com/greenlight/api/movie/MovieController.java`
- **Purpose:** Expose CRUD endpoints wired to `MovieService`.

```java
package com.greenlight.api.movie;

import com.greenlight.api.movie.dto.CreateMovieRequest;
import com.greenlight.api.movie.dto.MovieResponse;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/movies")
public class MovieController {

    private final MovieService movieService;

    public MovieController(MovieService movieService) {
        this.movieService = movieService;
    }

    @PostMapping
    public ResponseEntity<MovieResponse> createMovie(@Valid @RequestBody CreateMovieRequest request) {
        MovieResponse created = movieService.createMovie(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(created);
    }

    @GetMapping("/{id}")
    public ResponseEntity<MovieResponse> getMovie(@PathVariable Long id) {
        return ResponseEntity.ok(movieService.getMovieById(id));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteMovie(@PathVariable Long id) {
        movieService.deleteMovie(id);
        return ResponseEntity.noContent().build();
    }
}
```

### Git Milestone
```bash
git checkout -b feat/module-4-persistence
git add .
git commit -m "feat(movie): implement flyway migrations, jpa entity, repository, and service"
git checkout main && git merge feat/module-4-persistence
git push origin main
```

---

## Module 5: Optimistic Locking & Concurrency Control

### What you are migrating from Go:
- `internal/data/movies.go`:
  ```sql
  UPDATE movies SET ... version = version + 1
  WHERE id = $7 AND version = $8 RETURNING version
  ```
  If 0 rows were updated, Go returned `ErrEditConflict` (409 Conflict).

### Spring Boot Concepts to Master:
1. **`@Version` Annotation:** Hibernate automatically manages the version column and throws `OptimisticLockingFailureException` on stale updates.
2. **Handling 409 Conflict:** Mapping `OptimisticLockingFailureException` in `@RestControllerAdvice`.
3. **Partial Updates (`PATCH`):** Applying non-null fields to the entity.

---

### Step-by-Step Implementation Guide

#### Step 1: Create `UpdateMovieRequest` DTO
- **File:** `src/main/java/com/greenlight/api/movie/dto/UpdateMovieRequest.java`
- **Purpose:** Request payload for partial updates (all fields optional).

```java
package com.greenlight.api.movie.dto;

import com.greenlight.api.common.validator.UniqueElements;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.Positive;
import jakarta.validation.constraints.Size;
import org.hibernate.validator.constraints.URL;
import java.util.List;

public record UpdateMovieRequest(
    @Size(max = 500, message = "must not be more than 500 bytes long")
    String title,

    @Min(value = 1888, message = "must be greater than 1888")
    Integer year,

    @Positive(message = "must be a positive integer")
    Integer runtime,

    @Size(min = 1, max = 5, message = "must contain between 1 and 5 genres")
    @UniqueElements(message = "must not contain duplicate values")
    List<String> genres,

    @URL(message = "must be a valid URL")
    String imageUrl,

    @Size(min = 10, max = 1000, message = "must be between 10 and 1000 characters")
    String synopsis
) {}
```

---

#### Step 2: Implement `updateMovie` in `MovieService`
- **File:** `src/main/java/com/greenlight/api/movie/MovieService.java`
- **Purpose:** Applies partial updates; Hibernate automatically checks `@Version` on commit.

```java
// Inside MovieService.java:

@Transactional
public MovieResponse updateMovie(Long id, UpdateMovieRequest request) {
    Movie movie = movieRepository.findById(id)
        .orElseThrow(() -> new ResourceNotFoundException("Movie not found"));

    if (request.title() != null) movie.setTitle(request.title());
    if (request.year() != null) movie.setYear(request.year());
    if (request.runtime() != null) movie.setRuntime(request.runtime());
    if (request.genres() != null) movie.setGenres(request.genres());
    if (request.imageUrl() != null) movie.setImageUrl(request.imageUrl());
    if (request.synopsis() != null) movie.setSynopsis(request.synopsis());

    // When the transaction commits, Hibernate executes:
    // UPDATE movies SET ..., version = version + 1 WHERE id = ? AND version = ?
    return MovieResponse.fromEntity(movie);
}
```

---

#### Step 3: Handle `OptimisticLockingFailureException` in `GlobalExceptionHandler`
- **File:** `src/main/java/com/greenlight/api/common/error/GlobalExceptionHandler.java`

```java
// Add to GlobalExceptionHandler.java:

@ExceptionHandler(org.springframework.dao.OptimisticLockingFailureException.class)
public ResponseEntity<ErrorResponse> handleOptimisticLockingConflict(Exception ex) {
    String message = "unable to update the record due to an edit conflict, please try again";
    return ResponseEntity.status(HttpStatus.CONFLICT).body(new ErrorResponse(message));
}
```

---

#### Step 4: Add `@PatchMapping` to `MovieController`
- **File:** `src/main/java/com/greenlight/api/movie/MovieController.java`

```java
// Add to MovieController.java:

@PatchMapping("/{id}")
public ResponseEntity<MovieResponse> updateMovie(
        @PathVariable Long id,
        @Valid @RequestBody UpdateMovieRequest request) {
    return ResponseEntity.ok(movieService.updateMovie(id, request));
}
```

### Git Milestone
```bash
git checkout -b feat/module-5-optimistic-locking
git add .
git commit -m "feat(movie): implement patch updates and optimistic locking conflict handling"
git checkout main && git merge feat/module-5-optimistic-locking
git push origin main
```

---

## Module 6: Advanced Querying — Pagination, Dynamic Sorting & Search

### What you are migrating from Go:
- `internal/data/filters.go` & `movies.go` `GetAll()`:
  - Title filter.
  - Genres filter.
  - Sort support (`-year`, `title`, fallback tie-breaker `id ASC`).
  - Pagination metadata (`first_page`, `last_page`, `total_records`).

### Spring Boot Concepts to Master:
1. **Spring Data `Pageable` & `Page<T>`:** Standardized pagination and sorting.
2. **JPA Specifications (`Specification<Movie>`):** Composable `WHERE` predicates without raw SQL string concatenation.
3. **Custom Paged Envelope:** Matching Go's `metadata` output format.

---

### Step-by-Step Implementation Guide

#### Step 1: Define Pagination `Metadata` and `MovieListResponse` DTOs
- **File:** `src/main/java/com/greenlight/api/common/dto/Metadata.java`
- **Purpose:** Output format matching Go's pagination metadata.

```java
package com.greenlight.api.common.dto;

import org.springframework.data.domain.Page;

public record Metadata(
    int currentPage,
    int pageSize,
    int firstPage,
    int lastPage,
    long totalRecords
) {
    public static Metadata fromPage(Page<?> page) {
        return new Metadata(
            page.getNumber() + 1, // Spring is 0-indexed; client is 1-indexed
            page.getSize(),
            1,
            page.getTotalPages() == 0 ? 1 : page.getTotalPages(),
            page.getTotalElements()
        );
    }
}
```

- **File:** `src/main/java/com/greenlight/api/movie/dto/MovieListResponse.java`
```java
package com.greenlight.api.movie.dto;

import com.greenlight.api.common.dto.Metadata;
import java.util.List;

public record MovieListResponse(
    List<MovieResponse> movies,
    Metadata metadata
) {}
```

---

#### Step 2: Implement Dynamic Query Specifications
- **File:** `src/main/java/com/greenlight/api/movie/MovieSpecifications.java`
- **Purpose:** Builds dynamic predicates based on query parameters.

```java
package com.greenlight.api.movie;

import jakarta.persistence.criteria.Join;
import jakarta.persistence.criteria.Predicate;
import org.springframework.data.jpa.domain.Specification;

import java.util.ArrayList;
import java.util.List;

public class MovieSpecifications {

    public static Specification<Movie> filterBy(String title, List<String> genres) {
        return (root, query, cb) -> {
            List<Predicate> predicates = new ArrayList<>();

            if (title != null && !title.isBlank()) {
                predicates.add(cb.like(cb.lower(root.get("title")), "%" + title.toLowerCase() + "%"));
            }

            if (genres != null && !genres.isEmpty()) {
                Join<Movie, String> genreJoin = root.join("genres");
                predicates.add(genreJoin.in(genres));
            }

            query.distinct(true);
            return cb.and(predicates.toArray(new Predicate[0]));
        };
    }
}
```

---

#### Step 3: Implement `listMovies` in `MovieService`
- **File:** `src/main/java/com/greenlight/api/movie/MovieService.java`

```java
// Add to MovieService.java:

public MovieListResponse listMovies(String title, List<String> genres, Pageable pageable) {
    Specification<Movie> spec = MovieSpecifications.filterBy(title, genres);
    Page<Movie> page = movieRepository.findAll(spec, pageable);

    List<MovieResponse> movies = page.getContent().stream()
        .map(MovieResponse::fromEntity)
        .toList();

    return new MovieListResponse(movies, Metadata.fromPage(page));
}
```

---

#### Step 4: Expose `GET /v1/movies` in `MovieController`
- **File:** `src/main/java/com/greenlight/api/movie/MovieController.java`

```java
// Add to MovieController.java:

@GetMapping
public ResponseEntity<MovieListResponse> listMovies(
        @RequestParam(required = false) String title,
        @RequestParam(required = false) List<String> genres,
        @RequestParam(defaultValue = "1") int page,
        @RequestParam(defaultValue = "20") int pageSize,
        @RequestParam(defaultValue = "id") String sort) {

    // Translate sort string (e.g., "-year" -> descending by year)
    Sort sorting = sort.startsWith("-") 
        ? Sort.by(sort.substring(1)).descending() 
        : Sort.by(sort).ascending();

    // Secondary tie-breaker by id
    sorting = sorting.and(Sort.by("id").ascending());

    Pageable pageable = PageRequest.of(page - 1, pageSize, sorting);
    return ResponseEntity.ok(movieService.listMovies(title, genres, pageable));
}
```

### Git Milestone
```bash
git checkout -b feat/module-6-querying
git add .
git commit -m "feat(movie): implement dynamic specifications, pagination, and sorting"
git checkout main && git merge feat/module-6-querying
git push origin main
```

---

## Module 7: User Management & Password Security

### What you are migrating from Go:
- `migrations/000004_create_users_table.up.sql`
- `internal/data/users.go`: `User` struct, bcrypt password hashing, `Insert()`, `GetByEmail()`.
- `cmd/api/users.go`: `registerUserHandler` (`POST /v1/users`).

### Spring Boot Concepts to Master:
1. **Spring Security's `PasswordEncoder`:** Modern, secure password hashing using `BCryptPasswordEncoder`.
2. **Handling Unique Email Constraints:** Preventing duplicate user registration.
3. **Separating Domain Entity from API representation.**

---

### Step-by-Step Implementation Guide

#### Step 1: Database Migration for Users
- **File:** `src/main/resources/db/migration/V4__create_users_table.sql`

```sql
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);
```

---

#### Step 2: Configure `PasswordEncoder` Bean
- **File:** `src/main/java/com/greenlight/api/config/SecurityConfig.java`
- **Purpose:** Declares the BCrypt bean for dependency injection.

```java
package com.greenlight.api.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

@Configuration
public class SecurityConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(12); // Cost factor 12, matching Go
    }
}
```

---

#### Step 3: Implement `User` JPA Entity
- **File:** `src/main/java/com/greenlight/api/user/User.java`

```java
package com.greenlight.api.user;

import jakarta.persistence.*;
import org.hibernate.annotations.CreationTimestamp;
import java.time.Instant;

@Entity
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(name = "password_hash", nullable = false)
    private byte[] passwordHash;

    @Column(nullable = false)
    private boolean activated = false;

    @Version
    @Column(nullable = false)
    private Integer version = 1;

    public User() {}

    // Getters and Setters
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public Instant getCreatedAt() { return createdAt; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public String getEmail() { return email; }
    public void setEmail(String email) { this.email = email; }
    public byte[] getPasswordHash() { return passwordHash; }
    public void setPasswordHash(byte[] passwordHash) { this.passwordHash = passwordHash; }
    public boolean isActivated() { return activated; }
    public void setActivated(boolean activated) { this.activated = activated; }
    public Integer getVersion() { return version; }
    public void setVersion(Integer version) { this.version = version; }
}
```

---

#### Step 4: Implement `UserRepository`
- **File:** `src/main/java/com/greenlight/api/user/UserRepository.java`

```java
package com.greenlight.api.user;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);
    boolean existsByEmail(String email);
}
```

---

#### Step 5: Implement User Registration DTOs
- **File:** `src/main/java/com/greenlight/api/user/dto/RegisterUserRequest.java`
```java
package com.greenlight.api.user.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record RegisterUserRequest(
    @NotBlank(message = "must be provided")
    @Size(max = 500, message = "must not be more than 500 bytes long")
    String name,

    @NotBlank(message = "must be provided")
    @Email(message = "must be a valid email address")
    String email,

    @NotBlank(message = "must be provided")
    @Size(min = 8, message = "must not be less than 8 bytes long")
    String password
) {}
```

- **File:** `src/main/java/com/greenlight/api/user/dto/UserResponse.java`
```java
package com.greenlight.api.user.dto;

import com.greenlight.api.user.User;
import java.time.Instant;

public record UserResponse(
    Long id,
    Instant createdAt,
    String name,
    String email,
    boolean activated
) {
    public static UserResponse fromEntity(User user) {
        return new UserResponse(
            user.getId(),
            user.getCreatedAt(),
            user.getName(),
            user.getEmail(),
            user.isActivated()
        );
    }
}
```

---

#### Step 6: Create `DuplicateEmailException`
- **File:** `src/main/java/com/greenlight/api/user/DuplicateEmailException.java`
```java
package com.greenlight.api.user;

public class DuplicateEmailException extends RuntimeException {
    public DuplicateEmailException() {
        super("a user with this email address already exists");
    }
}
```
*(Add an `@ExceptionHandler` in `GlobalExceptionHandler` returning HTTP 422 with `{"email": "a user with this email address already exists"}`)*.

---

#### Step 7: Implement `UserService` and `UserController`
- **File:** `src/main/java/com/greenlight/api/user/UserService.java`

```java
package com.greenlight.api.user;

import com.greenlight.api.user.dto.RegisterUserRequest;
import com.greenlight.api.user.dto.UserResponse;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.nio.charset.StandardCharsets;

@Service
@Transactional(readOnly = true)
public class UserService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;

    public UserService(UserRepository userRepository, PasswordEncoder passwordEncoder) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
    }

    @Transactional
    public UserResponse registerUser(RegisterUserRequest request) {
        if (userRepository.existsByEmail(request.email())) {
            throw new DuplicateEmailException();
        }

        User user = new User();
        user.setName(request.name());
        user.setEmail(request.email());
        user.setPasswordHash(passwordEncoder.encode(request.password()).getBytes(StandardCharsets.UTF_8));
        user.setActivated(false);

        User saved = userRepository.save(user);
        return UserResponse.fromEntity(saved);
    }
}
```

- **File:** `src/main/java/com/greenlight/api/user/UserController.java`
```java
package com.greenlight.api.user;

import com.greenlight.api.user.dto.RegisterUserRequest;
import com.greenlight.api.user.dto.UserResponse;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @PostMapping
    public ResponseEntity<UserResponse> registerUser(@Valid @RequestBody RegisterUserRequest request) {
        UserResponse response = userService.registerUser(request);
        return ResponseEntity.status(HttpStatus.ACCEPTED).body(response);
    }
}
```

### Git Milestone
```bash
git checkout -b feat/module-7-users
git add .
git commit -m "feat(user): implement user registration and bcrypt password hashing"
git checkout main && git merge feat/module-7-users
git push origin main
```

---

## Module 8: Asynchronous Processing & Email Dispatch

### What you are migrating from Go:
- `cmd/api/server.go`: `app.background()` with `sync.WaitGroup` to run tasks asynchronously without blocking the HTTP response, while waiting for in-flight tasks during graceful shutdown.
- `internal/data/tokens.go`: Generating activation tokens (Base32 encoded random bytes, SHA-256 stored in DB).
- `internal/mailer/mailer.go`: Sending welcome email.
- `cmd/api/users.go` `activateUserHandler`: `PUT /v1/users/activated`.

### Spring Boot Concepts to Master:
1. **`@EnableAsync` and `@Async`:** Declarative background task execution.
2. **`ThreadPoolTaskExecutor`:** Configuring thread pools with graceful shutdown (Spring automatically waits for active async tasks to complete during shutdown, matching Go's `sync.WaitGroup`).
3. **Decoupled Architecture with Spring Events:** Publishing `UserRegisteredEvent` via `ApplicationEventPublisher`.

---

### Step-by-Step Implementation Guide

#### Step 1: Configure Asynchronous Execution
- **File:** `src/main/java/com/greenlight/api/config/AsyncConfig.java`
- **Purpose:** Thread pool with graceful shutdown.

```java
package com.greenlight.api.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableAsync;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import java.util.concurrent.Executor;

@Configuration
@EnableAsync
public class AsyncConfig {

    @Bean(name = "taskExecutor")
    public Executor taskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(500);
        executor.setThreadNamePrefix("Greenlight-Async-");
        executor.setWaitForTasksToCompleteOnShutdown(true); // Matches sync.WaitGroup
        executor.setAwaitTerminationSeconds(20);
        executor.initialize();
        return executor;
    }
}
```

---

#### Step 2: Database Migration for Tokens
- **File:** `src/main/resources/db/migration/V5__create_tokens_table.sql`

```sql
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
```

---

#### Step 3: Implement `Token` JPA Entity and `TokenService`
- **File:** `src/main/java/com/greenlight/api/token/Token.java`
```java
package com.greenlight.api.token;

import com.greenlight.api.user.User;
import jakarta.persistence.*;
import java.time.Instant;

@Entity
@Table(name = "tokens")
public class Token {

    @Id
    @Column(name = "hash")
    private byte[] hash;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Column(name = "expiry", nullable = false)
    private Instant expiry;

    @Column(name = "scope", nullable = false)
    private String scope;

    public Token() {}

    public byte[] getHash() { return hash; }
    public void setHash(byte[] hash) { this.hash = hash; }
    public User getUser() { return user; }
    public void setUser(User user) { this.user = user; }
    public Instant getExpiry() { return expiry; }
    public void setExpiry(Instant expiry) { this.expiry = expiry; }
    public String getScope() { return scope; }
    public void setScope(String scope) { this.scope = scope; }
}
```

- **File:** `src/main/java/com/greenlight/api/token/TokenRepository.java`
```java
package com.greenlight.api.token;

import com.greenlight.api.user.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.time.Instant;
import java.util.Optional;

@Repository
public interface TokenRepository extends JpaRepository<Token, byte[]> {
    Optional<Token> findByHashAndScopeAndExpiryAfter(byte[] hash, String scope, Instant now);
    void deleteAllByUserAndScope(User user, String scope);
}
```

- **File:** `src/main/java/com/greenlight/api/token/TokenService.java`
- **Purpose:** Generates cryptographically secure random tokens and stores SHA-256 hashes.

```java
package com.greenlight.api.token;

import com.greenlight.api.user.User;
import org.apache.commons.codec.binary.Base32;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.time.Duration;
import java.time.Instant;

@Service
public class TokenService {

    public static final String SCOPE_ACTIVATION = "activation";
    public static final String SCOPE_AUTHENTICATION = "authentication";

    private final TokenRepository tokenRepository;
    private final SecureRandom secureRandom = new SecureRandom();
    private final Base32 base32 = new Base32();

    public TokenService(TokenRepository tokenRepository) {
        this.tokenRepository = tokenRepository;
    }

    public record GeneratedToken(String plaintext, Token tokenEntity) {}

    @Transactional
    public GeneratedToken createToken(User user, Duration ttl, String scope) {
        byte[] randomBytes = new byte[16];
        secureRandom.nextBytes(randomBytes);
        String plaintext = base32.encodeAsString(randomBytes).replace("=", "");

        byte[] hash = sha256(plaintext);

        Token token = new Token();
        token.setHash(hash);
        token.setUser(user);
        token.setExpiry(Instant.now().plus(ttl));
        token.setScope(scope);

        Token saved = tokenRepository.save(token);
        return new GeneratedToken(plaintext, saved);
    }

    public byte[] sha256(String plaintext) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            return digest.digest(plaintext.getBytes());
        } catch (NoSuchAlgorithmException e) {
            throw new IllegalStateException("SHA-256 not available", e);
        }
    }
}
```

---

#### Step 4: Implement `MailService`
- **File:** `src/main/java/com/greenlight/api/mail/MailService.java`

```java
package com.greenlight.api.mail;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

@Service
public class MailService {

    private static final Logger log = LoggerFactory.getLogger(MailService.class);

    public void sendWelcomeEmail(String recipientEmail, String recipientName, String activationToken) {
        // In production: use JavaMailSender or Resend SDK
        log.info("Sending welcome email to {} ({}) with activation token: {}",
            recipientEmail, recipientName, activationToken);
    }
}
```

---

#### Step 5: Implement Decoupled User Registration Events
- **File:** `src/main/java/com/greenlight/api/user/event/UserRegisteredEvent.java`
```java
package com.greenlight.api.user.event;

import com.greenlight.api.user.User;

public record UserRegisteredEvent(User user, String plaintextToken) {}
```

- **File:** `src/main/java/com/greenlight/api/user/event/UserEventListener.java`
```java
package com.greenlight.api.user.event;

import com.greenlight.api.mail.MailService;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Component;
import org.springframework.transaction.event.TransactionPhase;
import org.springframework.transaction.event.TransactionalEventListener;

@Component
public class UserEventListener {

    private final MailService mailService;

    public UserEventListener(MailService mailService) {
        this.mailService = mailService;
    }

    @Async("taskExecutor")
    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void onUserRegistered(UserRegisteredEvent event) {
        mailService.sendWelcomeEmail(
            event.user().getEmail(),
            event.user().getName(),
            event.plaintextToken()
        );
    }
}
```

---

#### Step 6: Implement Account Activation (`PUT /v1/users/activated`)
- **File:** `src/main/java/com/greenlight/api/user/dto/ActivateUserRequest.java`
```java
package com.greenlight.api.user.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record ActivateUserRequest(
    @NotBlank(message = "must be provided")
    @Size(min = 26, max = 26, message = "must be 26 characters long")
    String token
) {}
```

In `UserService.java`, implement `activateUser(String plaintextToken)`:
- Hash token with `tokenService.sha256(plaintextToken)`.
- Find token in `tokenRepository.findByHashAndScopeAndExpiryAfter(...)`.
- Set user `activated = true`, save user, and call `tokenRepository.deleteAllByUserAndScope(user, SCOPE_ACTIVATION)`.

### Git Milestone
```bash
git checkout -b feat/module-8-async-mail
git add .
git commit -m "feat(user): implement token generation, async email event, and activation"
git checkout main && git merge feat/module-8-async-mail
git push origin main
```

---

## Module 9: Authentication & Stateful Bearer Tokens via Spring Security Filter Chain

### What you are migrating from Go:
- `cmd/api/tokens.go`: `POST /v1/tokens/authentication` creating a 24-hour token.
- `cmd/api/middleware.go` `authenticate()`: Extracts `Authorization: Bearer <token>`, looks up user in database, stores user in request context.

### Spring Boot Concepts to Master:
1. **Spring Security 6 Architecture:** `SecurityFilterChain`, `SecurityContextHolder`, `Authentication`.
2. **`OncePerRequestFilter`:** Custom filter to parse `Bearer <token>` on every request and populate the `SecurityContext`.
3. **Stateless Session Management:** Configuring `SessionCreationPolicy.STATELESS`.

---

### Step-by-Step Implementation Guide

#### Step 1: Authentication Token DTOs and Controller
- **File:** `src/main/java/com/greenlight/api/token/dto/CreateTokenRequest.java`
```java
package com.greenlight.api.token.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

public record CreateTokenRequest(
    @NotBlank(message = "must be provided")
    @Email(message = "must be a valid email address")
    String email,

    @NotBlank(message = "must be provided")
    String password
) {}
```

- **File:** `src/main/java/com/greenlight/api/token/dto/AuthenticationTokenResponse.java`
```java
package com.greenlight.api.token.dto;

import java.time.Instant;

public record AuthenticationTokenResponse(
    AuthenticationTokenDetails authentication_token
) {
    public record AuthenticationTokenDetails(String token, Instant expiry) {}
}
```

- **File:** `src/main/java/com/greenlight/api/token/TokenController.java`
```java
package com.greenlight.api.token;

import com.greenlight.api.token.dto.AuthenticationTokenResponse;
import com.greenlight.api.token.dto.CreateTokenRequest;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/tokens")
public class TokenController {

    private final TokenService tokenService;

    public TokenController(TokenService tokenService) {
        this.tokenService = tokenService;
    }

    @PostMapping("/authentication")
    public ResponseEntity<AuthenticationTokenResponse> createAuthToken(@Valid @RequestBody CreateTokenRequest request) {
        // TODO: Verify credentials and issue 24h authentication token
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }
}
```

---

#### Step 2: Implement `BearerTokenAuthenticationFilter`
- **File:** `src/main/java/com/greenlight/api/config/security/BearerTokenAuthenticationFilter.java`
- **Purpose:** Extracts Bearer token, queries DB, and populates `SecurityContext`.

```java
package com.greenlight.api.config.security;

import com.greenlight.api.token.TokenRepository;
import com.greenlight.api.token.TokenService;
import com.greenlight.api.user.User;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.HttpHeaders;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.time.Instant;
import java.util.Collections;

@Component
public class BearerTokenAuthenticationFilter extends OncePerRequestFilter {

    private final TokenService tokenService;
    private final TokenRepository tokenRepository;

    public BearerTokenAuthenticationFilter(TokenService tokenService, TokenRepository tokenRepository) {
        this.tokenService = tokenService;
        this.tokenRepository = tokenRepository;
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {

        response.addHeader("Vary", "Authorization");

        String authHeader = request.getHeader(HttpHeaders.AUTHORIZATION);
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String plaintextToken = authHeader.substring(7);
            byte[] hash = tokenService.sha256(plaintextToken);

            tokenRepository.findByHashAndScopeAndExpiryAfter(hash, TokenService.SCOPE_AUTHENTICATION, Instant.now())
                .ifPresent(token -> {
                    User user = token.getUser();
                    UsernamePasswordAuthenticationToken auth =
                        new UsernamePasswordAuthenticationToken(user, null, Collections.emptyList());
                    auth.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                    SecurityContextHolder.getContext().setAuthentication(auth);
                });
        }

        filterChain.doFilter(request, response);
    }
}
```

---

#### Step 3: Configure `SecurityFilterChain`
- **File:** `src/main/java/com/greenlight/api/config/SecurityConfig.java`

```java
// Update SecurityConfig.java:

@Bean
public SecurityFilterChain securityFilterChain(
        HttpSecurity http,
        BearerTokenAuthenticationFilter tokenFilter) throws Exception {

    return http
        .csrf(AbstractHttpConfigurer::disable)
        .sessionManagement(sm -> sm.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
        .authorizeHttpRequests(auth -> auth
            .requestMatchers("/v1/healthcheck").permitAll()
            .requestMatchers(HttpMethod.POST, "/v1/users").permitAll()
            .requestMatchers(HttpMethod.PUT, "/v1/users/activated").permitAll()
            .requestMatchers(HttpMethod.POST, "/v1/tokens/authentication").permitAll()
            .anyRequest().authenticated()
        )
        .addFilterBefore(tokenFilter, UsernamePasswordAuthenticationFilter.class)
        .exceptionHandling(ex -> ex
            .authenticationEntryPoint((req, res, e) -> {
                res.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                res.setHeader("WWW-Authenticate", "Bearer");
                res.setContentType("application/json");
                res.getWriter().write("{\"error\": \"you must be authenticated to access this resource\"}");
            })
        )
        .build();
}
```

### Git Milestone
```bash
git checkout -b feat/module-9-bearer-auth
git add .
git commit -m "feat(security): implement bearer token authentication filter and stateless chain"
git checkout main && git merge feat/module-9-bearer-auth
git push origin main
```

---

## Module 10: Role-Based Access Control (RBAC) & Method Security

### What you are migrating from Go:
- `migrations/000006_add_permissions.up.sql`: `permissions` and `users_permissions` tables.
- `internal/data/permissions.go`: `movies:read`, `movies:write`.
- `cmd/api/middleware.go`: `requireActivatedUser()` and `requirePermission("movies:read")`.

### Spring Boot Concepts to Master:
1. **Spring Method Security (`@EnableMethodSecurity`, `@PreAuthorize`):** Declarative authorization on controller methods.
2. **`GrantedAuthority`:** Mapping permission codes into Spring Security authorities.
3. **SpEL (Spring Expression Language):** Checking permissions and activation status.

---

### Step-by-Step Implementation Guide

#### Step 1: Database Migration for Permissions
- **File:** `src/main/resources/db/migration/V6__add_permissions.sql`

```sql
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code) VALUES ('movies:read'), ('movies:write');
```

---

#### Step 2: Implement `Permission` JPA Entity and Map to `User`
- **File:** `src/main/java/com/greenlight/api/user/Permission.java`
```java
package com.greenlight.api.user;

import jakarta.persistence.*;

@Entity
@Table(name = "permissions")
public class Permission {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, unique = true)
    private String code;

    public Permission() {}

    public Long getId() { return id; }
    public String getCode() { return code; }
    public void setCode(String code) { this.code = code; }
}
```

In `User.java`, add the `@ManyToMany` mapping:
```java
// Inside User.java:

@ManyToMany(fetch = FetchType.EAGER)
@JoinTable(
    name = "users_permissions",
    joinColumns = @JoinColumn(name = "user_id"),
    inverseJoinColumns = @JoinColumn(name = "permission_id")
)
private Set<Permission> permissions = new HashSet<>();

public Set<Permission> getPermissions() { return permissions; }
```

---

#### Step 3: Enable Method Security
- **File:** `src/main/java/com/greenlight/api/config/SecurityConfig.java`

```java
// Add to SecurityConfig class:
@Configuration
@EnableMethodSecurity(prePostEnabled = true)
public class SecurityConfig { ... }
```

In `BearerTokenAuthenticationFilter.java`, map user permissions to authorities:
```java
List<SimpleGrantedAuthority> authorities = user.getPermissions().stream()
    .map(p -> new SimpleGrantedAuthority(p.getCode()))
    .toList();

UsernamePasswordAuthenticationToken auth =
    new UsernamePasswordAuthenticationToken(user, null, authorities);
```

---

#### Step 4: Protect Endpoints in `MovieController`
- **File:** `src/main/java/com/greenlight/api/movie/MovieController.java`

```java
// Secure methods in MovieController.java:

@GetMapping
@PreAuthorize("hasAuthority('movies:read')")
public ResponseEntity<MovieListResponse> listMovies(...) { ... }

@PostMapping
@PreAuthorize("hasAuthority('movies:write') and principal.activated")
public ResponseEntity<MovieResponse> createMovie(...) { ... }

@PatchMapping("/{id}")
@PreAuthorize("hasAuthority('movies:write') and principal.activated")
public ResponseEntity<MovieResponse> updateMovie(...) { ... }

@DeleteMapping("/{id}")
@PreAuthorize("hasAuthority('movies:write') and principal.activated")
public ResponseEntity<Void> deleteMovie(...) { ... }
```

---

#### Step 5: Configure 403 Forbidden Handling
In `SecurityConfig.java`, configure `accessDeniedHandler`:
```java
.exceptionHandling(ex -> ex
    .accessDeniedHandler((req, res, e) -> {
        res.setStatus(HttpServletResponse.SC_FORBIDDEN);
        res.setContentType("application/json");
        res.getWriter().write("{\"error\": \"your user account doesn't have the necessary permissions to access this resource\"}");
    })
)
```

### Git Milestone
```bash
git checkout -b feat/module-10-rbac
git add .
git commit -m "feat(security): implement rbac permissions and method security annotations"
git checkout main && git merge feat/module-10-rbac
git push origin main
```

---

## Module 11: Rate Limiting & CORS Policies

### What you are migrating from Go:
- `cmd/api/middleware.go` `rateLimit()`: Token bucket per IP (2 requests/sec, burst 4).
- `cmd/api/middleware.go` `enableCORS()`: Handling `Origin`, preflight `OPTIONS`, allowed headers.

### Spring Boot Concepts to Master:
1. **Token Bucket Rate Limiting with Bucket4j:** Standard production library for rate limiting in Java.
2. **Spring WebMvc CORS Configuration:** Using `CorsConfigurationSource` to manage trusted origins, allowed methods, and preflight requests without manual middleware closures.

---

### Step-by-Step Implementation Guide

#### Step 1: Add Bucket4j Dependency
- **File:** `pom.xml`

```xml
<dependency>
    <groupId>com.bucket4j</groupId>
    <artifactId>bucket4j-core</artifactId>
    <version>8.10.1</version>
</dependency>
```

---

#### Step 2: Implement `RateLimitingFilter`
- **File:** `src/main/java/com/greenlight/api/config/filter/RateLimitingFilter.java`
- **Purpose:** In-memory token-bucket limiter per client IP.

```java
package com.greenlight.api.config.filter;

import io.github.bucket4j.Bandwidth;
import io.github.bucket4j.Bucket;
import io.github.bucket4j.Refill;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
public class RateLimitingFilter extends OncePerRequestFilter {

    private final Map<String, Bucket> buckets = new ConcurrentHashMap<>();

    private Bucket createNewBucket() {
        // 2 requests per second with burst of 4 (matches Go limiter)
        Refill refill = Refill.greedy(2, Duration.ofSeconds(1));
        Bandwidth limit = Bandwidth.classic(4, refill);
        return Bucket.builder().addLimit(limit).build();
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {

        String ip = request.getRemoteAddr();
        Bucket bucket = buckets.computeIfAbsent(ip, k -> createNewBucket());

        if (bucket.tryConsume(1)) {
            filterChain.doFilter(request, response);
        } else {
            response.setStatus(429);
            response.setContentType("application/json");
            response.getWriter().write("{\"error\": \"rate limit exceeded\"}");
        }
    }
}
```

---

#### Step 3: Configure CORS in `SecurityConfig`
- **File:** `src/main/java/com/greenlight/api/config/SecurityConfig.java`

```java
// Inside SecurityConfig.java:

@Bean
public CorsConfigurationSource corsConfigurationSource(
        @Value("${app.cors.trusted-origins:}") List<String> trustedOrigins) {
    CorsConfiguration configuration = new CorsConfiguration();
    configuration.setAllowedOrigins(trustedOrigins);
    configuration.setAllowedMethods(List.of("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"));
    configuration.setAllowedHeaders(List.of("Authorization", "Content-Type"));
    configuration.setMaxAge(3600L);

    UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
    source.registerCorsConfiguration("/**", configuration);
    return source;
}
```
In your `securityFilterChain`:
```java
http.cors(cors -> cors.configurationSource(corsConfigurationSource(trustedOrigins)));
```

### Git Milestone
```bash
git checkout -b feat/module-11-rate-limit-cors
git add .
git commit -m "feat(web): configure bucket4j rate limiting filter and cors source"
git checkout main && git merge feat/module-11-rate-limit-cors
git push origin main
```

---

## Module 12: Production Observability, Metrics & Actuator

### What you are migrating from Go:
- `cmd/api/middleware.go` `metrics()`: Counting requests, response statuses, total processing time.
- `cmd/api/main.go` `expvar`: Goroutine counts, DB pool stats at `/debug/vars`.

### Spring Boot Concepts to Master:
1. **Spring Boot Actuator:** Production-ready endpoints for health, metrics, environment, threads.
2. **Micrometer:** Application metrics facade tracking JVM memory, active HikariCP connections, HTTP request durations and status codes.

---

### Step-by-Step Implementation Guide

#### Step 1: Configure Actuator in `application.yml`
- **File:** `src/main/resources/application.yml`

```yaml
management:
  endpoints:
    web:
      exposure:
        include: "health,info,metrics,prometheus"
  endpoint:
    health:
      show-details: when_authorized
  metrics:
    tags:
      application: ${spring.application.name}
```

---

#### Step 2: Permit Actuator in `SecurityConfig`
- **File:** `src/main/java/com/greenlight/api/config/SecurityConfig.java`

```java
// Add to securityFilterChain requestMatchers:
.requestMatchers("/actuator/**").permitAll()
```

### Verification
Execute in terminal:
- `curl http://localhost:4000/actuator/health`
- `curl http://localhost:4000/actuator/metrics`
- `curl http://localhost:4000/actuator/metrics/http.server.requests` (replaces Go's custom `httpsnoop` middleware!)
- `curl http://localhost:4000/actuator/metrics/hikaricp.connections.active` (replaces Go's `db.Stats()`)

### Git Milestone
```bash
git checkout -b feat/module-12-observability
git add .
git commit -m "feat(observability): configure spring boot actuator and micrometer metrics"
git checkout main && git merge feat/module-12-observability
git push origin main
```

---

## Chapter 16: Masterclass: Comprehensive Testing Strategy

In Chapter 20 of *Let's Go Further*, you wrote tests with `httptest.NewServer`.
In Spring Boot, testing spans across the **Testing Pyramid**:

```
       ▲
      / \     Tier 4: End-to-End (@SpringBootTest + Testcontainers)
     /   \
    /     \   Tier 2 & 3: Slices (@WebMvcTest, @DataJpaTest)
   /       \
  /_________\ Tier 1: Pure Unit Tests (JUnit 5 + Mockito, sub-millisecond)
```

---

### Tier 1: Pure Unit Testing with Mockito
- **File:** `src/test/java/com/greenlight/api/movie/MovieServiceTest.java`
- **Purpose:** Test isolated business logic in `MovieService` without booting any Spring context.

```java
package com.greenlight.api.movie;

import com.greenlight.api.common.error.ResourceNotFoundException;
import com.greenlight.api.movie.dto.MovieResponse;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class MovieServiceTest {

    @Mock
    private MovieRepository movieRepository;

    @InjectMocks
    private MovieService movieService;

    @Test
    void getMovieById_WhenMovieExists_ShouldReturnMovieResponse() {
        Movie movie = new Movie();
        movie.setId(1L);
        movie.setTitle("Inception");
        when(movieRepository.findById(1L)).thenReturn(Optional.of(movie));

        MovieResponse response = movieService.getMovieById(1L);

        assertNotNull(response);
        assertEquals("Inception", response.title());
        verify(movieRepository, times(1)).findById(1L);
    }

    @Test
    void getMovieById_WhenNotFound_ShouldThrowResourceNotFoundException() {
        when(movieRepository.findById(99L)).thenReturn(Optional.empty());

        assertThrows(ResourceNotFoundException.class, () -> movieService.getMovieById(99L));
    }
}
```

---

### Tier 2: Web Layer Slice Testing (`@WebMvcTest`)
- **File:** `src/test/java/com/greenlight/api/movie/MovieControllerTest.java`
- **Purpose:** Fast tests (<200ms) for validation rules, routing, status codes, and JSON serialization.

```java
package com.greenlight.api.movie;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@WebMvcTest(MovieController.class)
@AutoConfigureMockMvc(addFilters = false) // Disable security filters for pure web layer test
class MovieControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private MovieService movieService;

    @Test
    void createMovie_WithInvalidPayload_ShouldReturn422UnprocessableEntity() throws Exception {
        String invalidJson = """
            {
                "title": "",
                "year": 1800,
                "runtime": -10,
                "genres": []
            }
            """;

        mockMvc.perform(post("/v1/movies")
                .contentType(MediaType.APPLICATION_JSON)
                .content(invalidJson))
            .andExpect(status().isUnprocessableEntity())
            .andExpect(jsonPath("$.error.title").value("must be provided"))
            .andExpect(jsonPath("$.error.year").value("must be greater than 1888"))
            .andExpect(jsonPath("$.error.runtime").value("must be a positive integer"))
            .andExpect(jsonPath("$.error.genres").exists());
    }
}
```

---

### Tier 3: Data Layer Slice Testing (`@DataJpaTest`)
- **File:** `src/test/java/com/greenlight/api/movie/MovieRepositoryTest.java`
- **Purpose:** Test JPA entity mappings, database constraints, and optimistic locking.

```java
package com.greenlight.api.movie;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.autoconfigure.orm.jpa.TestEntityManager;
import org.springframework.dao.OptimisticLockingFailureException;

import static org.junit.jupiter.api.Assertions.assertThrows;

@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class MovieRepositoryTest {

    @Autowired
    private MovieRepository movieRepository;

    @Autowired
    private TestEntityManager entityManager;

    @Test
    void updateMovie_WithConflictingVersion_ShouldThrowOptimisticLockException() {
        Movie movie = new Movie();
        movie.setTitle("Original Title");
        movie.setYear(2020);
        movie.setRuntime(100);
        movie.setImageUrl("https://example.com/img.jpg");
        movie.setSynopsis("A great movie synopsis here.");
        Movie saved = movieRepository.saveAndFlush(movie);

        saved.setTitle("Updated Title 1");
        movieRepository.saveAndFlush(saved);

        // Force stale update with old version
        entityManager.detach(saved);
        saved.setVersion(0);
        saved.setTitle("Updated Title 2");

        assertThrows(OptimisticLockingFailureException.class, () -> {
            movieRepository.saveAndFlush(saved);
        });
    }
}
```

---

### Tier 4: End-to-End Integration Testing with Testcontainers
- **File:** `src/test/java/com/greenlight/api/GreenlightE2EIntegrationTest.java`
- **Purpose:** Boot the entire Spring Boot application against a real PostgreSQL container.

```java
package com.greenlight.api;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
@Testcontainers
class GreenlightE2EIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private MockMvc mockMvc;

    @Test
    void healthCheck_ShouldReturnAvailable() throws Exception {
        mockMvc.perform(get("/v1/healthcheck"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.status").value("available"));
    }

    @Test
    void protectedMoviesEndpoint_WithoutAuthentication_ShouldReturn401() throws Exception {
        mockMvc.perform(get("/v1/movies"))
            .andExpect(status().isUnauthorized());
    }
}
```

---

## Appendix: Go to Spring Boot Rosetta Stone

| Concept in Go (`greenlight`) | Idiomatic Spring Boot 3 Equivalent |
| :--- | :--- |
| `type application struct { ... }` | Spring IoC Container (`ApplicationContext`) |
| `app.models.Movies.Insert(&movie)` | `movieRepository.save(movie)` via Spring Data JPA |
| `app.writeJSON(w, status, data, headers)` | `ResponseEntity.status(status).headers(h).body(data)` |
| `v := validator.New(); v.Check(...)` | Jakarta Validation annotations (`@Valid`, `@NotBlank`, etc.) |
| `app.serverErrorResponse()` / `errors.go` | `@RestControllerAdvice` + `@ExceptionHandler` |
| `WHERE id = $1 AND version = $2` | `@Version` field on JPA entity |
| Goroutine `go func() { ... }()` with `sync.WaitGroup` | `@Async` + `ThreadPoolTaskExecutor.setWaitForTasksToCompleteOnShutdown(true)` |
| `app.metrics(...)` / `app.recoverPanic(...)` | Servlet Filter Chain & Spring Security Filter Chain |
| `to_tsvector` and dynamic query builder | `JpaSpecificationExecutor<Movie>` + `Specification<Movie>` |
| `bcrypt.GenerateFromPassword(pwd, 12)` | `new BCryptPasswordEncoder(12).encode(pwd)` |
| `r = app.contextSetUser(r, user)` | `SecurityContextHolder.getContext().setAuthentication(auth)` |
| `requirePermission("movies:write")` | `@PreAuthorize("hasAuthority('movies:write')")` |
| `rate.NewLimiter(...)` per IP | Bucket4j token bucket filter |
| `expvar` `/debug/vars` | Spring Boot Actuator (`/actuator/metrics`, `/actuator/health`) |

---

## Next Steps in IntelliJ IDEA

1. Create a new Spring Boot project using **Module 0**.
2. Follow each module sequentially. Create the specified file at the given path, implement the blueprint skeleton, run the verification command, and push the branch to GitHub.
3. Finish by executing the **Testing Masterclass** suite in IntelliJ IDEA!
