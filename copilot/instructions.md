Refactor this entire repository layer to properly support transaction propagation using the WithTx pattern.

Requirements:

1. Every repository struct that contains `DB *gorm.DB` must implement:

   func (r *RepoName) WithTx(tx *gorm.DB) *RepoName {
       return &RepoName{DB: tx}
   }

2. All repository methods must use `r.DB` internally and must not use any global DB reference.

3. In SharedRepo methods where a transaction is started:

   tx := s.DB.Begin()

   Replace direct usage of s.UserRepo, s.RoleRepo, s.ResetCredsRepo, etc with:

   userRepo := s.UserRepo.WithTx(tx)
   roleRepo := s.RoleRepo.WithTx(tx)
   resetRepo := s.ResetCredsRepo.WithTx(tx)

   Then use those cloned repos inside the transaction.

4. Ensure all database operations inside SharedRepo transactional methods use the transaction-bound repos.

5. Do NOT change business logic, method signatures, or return types.

6. Ensure rollback happens via defer tx.Rollback() and commit only on success.

7. Preserve existing error handling.

Apply this consistently across:
- UserRepository
- ResetCredsRepository
- RoleRepository
- RouteRoleRepository
- TokenRepository
- Any repository that contains DB *gorm.DB

Make the changes minimal and safe.