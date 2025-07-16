import React from 'react';

const AuthProLanding = () => {
  return (
    <div className="relative flex size-full min-h-screen flex-col bg-[#101623] overflow-x-hidden font-[Inter,'Noto Sans',sans-serif]">
      <div className="layout-container flex h-full grow flex-col">
        {/* Header */}
        <header className="flex items-center justify-between whitespace-nowrap border-b border-solid border-b-[#222f49] px-4 py-3 md:px-10">
          <div className="flex items-center gap-4 text-white">
            <div className="size-4">
              <LogoIcon />
            </div>
            <h2 className="text-white text-lg font-bold leading-tight tracking-[-0.015em]">AuthPro</h2>
          </div>
          <div className="flex flex-1 justify-end gap-4 md:gap-8">
            <div className="hidden items-center gap-5 md:flex md:gap-9">
              <NavLink href="#">Features</NavLink>
              <NavLink href="#">Pricing</NavLink>
              <NavLink href="#">Documentation</NavLink>
              <NavLink href="#">Support</NavLink>
            </div>
            <PrimaryButton className="hidden md:flex">Sign Up</PrimaryButton>
            <button className="md:hidden text-white">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
            </button>
          </div>
        </header>

        {/* Main Content */}
        <main className="px-4 flex flex-1 justify-center py-5 md:px-10 lg:px-40">
          <div className="layout-content-container flex flex-col max-w-[960px] flex-1">
            {/* Hero Section */}
            <div className="@container">
              <div className="flex flex-col gap-6 px-0 py-10 @[480px]:gap-8 @[864px]:flex-row">
                <div 
                  className="w-full bg-center bg-no-repeat aspect-video bg-cover rounded-lg @[480px]:h-auto @[480px]:min-w-[400px] @[864px]:w-full"
                  style={{ backgroundImage: 'url(https://lh3.googleusercontent.com/aida-public/AB6AXuB8MGAYp-LLtm4vxCkcr8EMRbsscaVh8PvpWZb9V-mkvV2AXgcgZDPeR4Au6RY5bRcuixh_ruzhtBCd5_klyzsC6IgveCj5dgAnI96ckU1KQRmxyavqLfAjgOofOTxriO0PpzAE-4gfDSaawYNXso-L3pWtbaLfba_QNPAC0KYCymDEdL_CRmZIN1-DNB_QgXtG9UDHAXMm30Hkxa1JP2TcZ9KDhxbeM4-bPdYjmpL6fGMTdGn4-zYMfzFRtdbcpva3fy3YOmJWeVM)' }}
                />
                <div className="flex flex-col gap-6 @[480px]:min-w-[400px] @[480px]:gap-8 @[864px]:justify-center">
                  <div className="flex flex-col gap-2 text-left">
                    <h1 className="text-white text-3xl font-black leading-tight tracking-[-0.033em] @[480px]:text-4xl @[480px]:text-5xl">
                      Plug-and-Play Authentication for Your Apps
                    </h1>
                    <h2 className="text-white text-sm font-normal leading-normal @[480px]:text-base">
                      Simplify user authentication with AuthPro's easy-to-integrate solution. Secure your applications quickly and efficiently, allowing you to focus on your core product.
                    </h2>
                  </div>
                  <PrimaryButton className="@[480px]:h-12 @[480px]:px-5 @[480px]:text-base">
                    Get Started
                  </PrimaryButton>
                </div>
              </div>
            </div>

            {/* Key Features Section */}
            <SectionTitle>Key Features</SectionTitle>
            <div className="flex flex-col gap-10 px-0 py-10 @container">
              <div className="flex flex-col gap-4">
                <h1 className="text-white tracking-light text-2xl font-bold leading-tight @[480px]:text-[32px] @[480px]:text-4xl @[480px]:font-black @[480px]:tracking-[-0.033em] max-w-[720px]">
                  Secure and Scalable Authentication
                </h1>
                <p className="text-white text-base font-normal leading-normal max-w-[720px]">
                  AuthPro provides a robust and scalable authentication solution, ensuring the security of your applications and user data.
                </p>
              </div>
              <div className="grid grid-cols-1 gap-3 p-0 sm:grid-cols-2 lg:grid-cols-4">
                <FeatureCard 
                  icon={<ShieldIcon />}
                  title="Multi-Factor Authentication"
                  description="Enhance security with multi-factor authentication, requiring users to provide multiple forms of verification."
                />
                <FeatureCard 
                  icon={<KeyIcon />}
                  title="Passwordless Login"
                  description="Enable passwordless login options, such as magic links or biometric authentication, for a seamless user experience."
                />
                <FeatureCard 
                  icon={<UsersIcon />}
                  title="User Management"
                  description="Manage user accounts, roles, and permissions with ease, providing granular control over access."
                />
                <FeatureCard 
                  icon={<LockIcon />}
                  title="Data Encryption"
                  description="Protect sensitive data with industry-standard encryption techniques, ensuring confidentiality and integrity."
                />
              </div>
            </div>

            {/* Customer Success Section */}
            <SectionTitle>Customer Success</SectionTitle>
            <div className="py-4">
              <div className="flex flex-col items-stretch justify-between gap-6 rounded-lg lg:flex-row">
                <div className="flex flex-[2_2_0px] flex-col gap-4">
                  <div className="flex flex-col gap-1">
                    <p className="text-[#90a4cb] text-sm font-normal leading-normal">Case Study</p>
                    <p className="text-white text-base font-bold leading-tight">Streamlining User Onboarding for QuickStart</p>
                    <p className="text-[#90a4cb] text-sm font-normal leading-normal">
                      QuickStart, a leading productivity app, integrated AuthPro to simplify user onboarding and improve security. The result was a 30% reduction in onboarding time
                      and a significant increase in user satisfaction.
                    </p>
                  </div>
                  <SecondaryButton>
                    Read More
                  </SecondaryButton>
                </div>
                <div 
                  className="w-full bg-center bg-no-repeat aspect-video bg-cover rounded-lg flex-1"
                  style={{ backgroundImage: 'url(https://lh3.googleusercontent.com/aida-public/AB6AXuCKUgSVp6q-t2hrZqjDN16KTR7I3XZ-adVFERySHei31h5TnXr0A7HaoBdrZUbIO-7C0yYgvRcdvS0LFCcmQ9OLcNBG9YuhBZQ1Q9QuVheO0fu7eZu9QliFzMYrAEZGvzAuKSZoI5aEOC5kQ6H1D9Bmc8mf5S1iGAVyWztobSJ34NU35lR6CN_fE38P21O7TUrF44jC45jZmbU1hpxhDGCoSAbxY5X4mJV7WIwNLFU5OgbxABHjA0dLcE87bhiFjYg6CBnRCGI37OQ)' }}
                />
              </div>
            </div>

            {/* CTA Section */}
            <div className="@container">
              <div className="flex flex-col justify-end gap-6 px-0 py-10 @[480px]:gap-8 @[480px]:px-10 @[480px]:py-20">
                <div className="flex flex-col gap-2 text-center">
                  <h1 className="text-white tracking-light text-2xl font-bold leading-tight @[480px]:text-[32px] @[480px]:text-4xl @[480px]:font-black @[480px]:tracking-[-0.033em] max-w-[720px] mx-auto">
                    Ready to Secure Your App?
                  </h1>
                </div>
                <div className="flex flex-1 justify-center">
                  <div className="flex justify-center w-full max-w-[480px]">
                    <PrimaryButton className="@[480px]:h-12 @[480px]:px-5 @[480px]:text-base w-full">
                      Get Started Today
                    </PrimaryButton>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>

        {/* Footer */}
        <footer className="flex justify-center">
          <div className="flex max-w-[960px] flex-1 flex-col">
            <div className="flex flex-col gap-6 px-5 py-10 text-center @container">
              <div className="flex flex-wrap items-center justify-center gap-6 @[480px]:flex-row @[480px]:justify-around">
                <FooterLink href="#">About Us</FooterLink>
                <FooterLink href="#">Terms of Service</FooterLink>
                <FooterLink href="#">Privacy Policy</FooterLink>
              </div>
              <div className="flex flex-wrap justify-center gap-4">
                <SocialIcon href="#">
                  <TwitterIcon />
                </SocialIcon>
                <SocialIcon href="#">
                  <LinkedInIcon />
                </SocialIcon>
              </div>
              <p className="text-[#90a4cb] text-base font-normal leading-normal">© 2024 AuthPro. All rights reserved.</p>
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
};

// Reusable Components
const NavLink = ({ href, children }) => (
  <a className="text-white text-sm font-medium leading-normal hover:text-[#2469f3] transition-colors" href={href}>
    {children}
  </a>
);

const PrimaryButton = ({ children, className = '' }) => (
  <button
    className={`flex min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-10 px-4 bg-[#2469f3] text-white text-sm font-bold leading-normal tracking-[0.015em] hover:bg-[#1a5bd9] transition-colors ${className}`}
  >
    <span className="truncate">{children}</span>
  </button>
);

const SecondaryButton = ({ children }) => (
  <button
    className="flex min-w-[84px] max-w-[480px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-8 px-4 flex-row-reverse bg-[#222f49] text-white text-sm font-medium leading-normal w-fit hover:bg-[#314368] transition-colors"
  >
    <span className="truncate">{children}</span>
  </button>
);

const SectionTitle = ({ children }) => (
  <h2 className="text-white text-[22px] font-bold leading-tight tracking-[-0.015em] px-0 pb-3 pt-5">
    {children}
  </h2>
);

const FeatureCard = ({ icon, title, description }) => (
  <div className="flex flex-1 gap-3 rounded-lg border border-[#314368] bg-[#182234] p-4 flex-col hover:border-[#2469f3] transition-colors">
    <div className="text-white">{icon}</div>
    <div className="flex flex-col gap-1">
      <h2 className="text-white text-base font-bold leading-tight">{title}</h2>
      <p className="text-[#90a4cb] text-sm font-normal leading-normal">{description}</p>
    </div>
  </div>
);

const FooterLink = ({ href, children }) => (
  <a className="text-[#90a4cb] text-base font-normal leading-normal hover:text-white transition-colors min-w-40" href={href}>
    {children}
  </a>
);

const SocialIcon = ({ href, children }) => (
  <a href={href} className="text-[#90a4cb] hover:text-white transition-colors">
    {children}
  </a>
);

// Icons
const LogoIcon = () => (
  <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M12.0799 24L4 19.2479L9.95537 8.75216L18.04 13.4961L18.0446 4H29.9554L29.96 13.4961L38.0446 8.75216L44 19.2479L35.92 24L44 28.7521L38.0446 39.2479L29.96 34.5039L29.9554 44H18.0446L18.04 34.5039L9.95537 39.2479L4 28.7521L12.0799 24Z"
      fill="currentColor"
    />
  </svg>
);

const ShieldIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M208,40H48A16,16,0,0,0,32,56v58.77c0,89.61,75.82,119.34,91,124.39a15.53,15.53,0,0,0,10,0c15.2-5.05,91-34.78,91-124.39V56A16,16,0,0,0,208,40Zm0,74.79c0,78.42-66.35,104.62-80,109.18-13.53-4.51-80-30.69-80-109.18V56l160,0Z" />
  </svg>
);

const KeyIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M160,16A80.07,80.07,0,0,0,83.91,120.78L26.34,178.34A8,8,0,0,0,24,184v40a8,8,0,0,0,8,8H72a8,8,0,0,0,8-8V208H96a8,8,0,0,0,8-8V184h16a8,8,0,0,0,5.66-2.34l9.56-9.57A80,80,0,1,0,160,16Zm0,144a63.7,63.7,0,0,1-23.65-4.51,8,8,0,0,0-8.84,1.68L116.69,168H96a8,8,0,0,0-8,8v16H72a8,8,0,0,0-8,8v16H40V187.31l58.83-58.82a8,8,0,0,0,1.68-8.84A64,64,0,1,1,160,160Zm32-84a12,12,0,1,1-12-12A12,12,0,0,1,192,76Z" />
  </svg>
);

const UsersIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M117.25,157.92a60,60,0,1,0-66.5,0A95.83,95.83,0,0,0,3.53,195.63a8,8,0,1,0,13.4,8.74,80,80,0,0,1,134.14,0,8,8,0,0,0,13.4-8.74A95.83,95.83,0,0,0,117.25,157.92ZM40,108a44,44,0,1,1,44,44A44.05,44.05,0,0,1,40,108Zm210.14,98.7a8,8,0,0,1-11.07-2.33A79.83,79.83,0,0,0,172,168a8,8,0,0,1,0-16,44,44,0,1,0-16.34-84.87,8,8,0,1,1-5.94-14.85,60,60,0,0,1,55.53,105.64,95.83,95.83,0,0,1,47.22,37.71A8,8,0,0,1,250.14,206.7Z" />
  </svg>
);

const LockIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M208,80H176V56a48,48,0,0,0-96,0V80H48A16,16,0,0,0,32,96V208a16,16,0,0,0,16,16H208a16,16,0,0,0,16-16V96A16,16,0,0,0,208,80ZM96,56a32,32,0,0,1,64,0V80H96ZM208,208H48V96H208V208Zm-68-56a12,12,0,1,1-12-12A12,12,0,0,1,140,152Z" />
  </svg>
);

const TwitterIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M247.39,68.94A8,8,0,0,0,240,64H209.57A48.66,48.66,0,0,0,168.1,40a46.91,46.91,0,0,0-33.75,13.7A47.9,47.9,0,0,0,120,88v6.09C79.74,83.47,46.81,50.72,46.46,50.37a8,8,0,0,0-13.65,4.92c-4.31,47.79,9.57,79.77,22,98.18a110.93,110.93,0,0,0,21.88,24.2c-15.23,17.53-39.21,26.74-39.47,26.84a8,8,0,0,0-3.85,11.93c.75,1.12,3.75,5.05,11.08,8.72C53.51,229.7,65.48,232,80,232c70.67,0,129.72-54.42,135.75-124.44l29.91-29.9A8,8,0,0,0,247.39,68.94Zm-45,29.41a8,8,0,0,0-2.32,5.14C196,166.58,143.28,216,80,216c-10.56,0-18-1.4-23.22-3.08,11.51-6.25,27.56-17,37.88-32.48A8,8,0,0,0,92,169.08c-.47-.27-43.91-26.34-44-96,16,13,45.25,33.17,78.67,38.79A8,8,0,0,0,136,104V88a32,32,0,0,1,9.6-22.92A30.94,30.94,0,0,1,167.9,56c12.66.16,24.49,7.88,29.44,19.21A8,8,0,0,0,204.67,80h16Z" />
  </svg>
);

const LinkedInIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="24px" height="24px" fill="currentColor" viewBox="0 0 256 256">
    <path d="M216,24H40A16,16,0,0,0,24,40V216a16,16,0,0,0,16,16H216a16,16,0,0,0,16-16V40A16,16,0,0,0,216,24Zm0,192H40V40H216V216ZM96,112v64a8,8,0,0,1-16,0V112a8,8,0,0,1,16,0Zm88,28v36a8,8,0,0,1-16,0V140a20,20,0,0,0-40,0v36a8,8,0,0,1-16,0V112a8,8,0,0,1,15.79-1.78A36,36,0,0,1,184,140ZM100,84A12,12,0,1,1,88,72,12,12,0,0,1,100,84Z" />
  </svg>
);

export default AuthProLanding;